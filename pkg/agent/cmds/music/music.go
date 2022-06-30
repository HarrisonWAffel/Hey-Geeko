package music

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"harrisonwaffel/assistant/pkg/conversation"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init() {
	// read in a sample file to get the format for all songs
	f, err := os.Open("pkg/agent/cmds/music/files/puppy_love.mp3")
	if err != nil {
		panic("could not initialize music app: " + err.Error())
	}
	_, format, err := mp3.Decode(f)
	if err != nil {
		panic("could not decode mp3 during music client initialization")
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	wd, _ := os.Getwd()
	musicFilePath = wd + "/pkg/cmds/music/files/"
}

var (
	// genreURL points to deezer for free audio previews
	// of various genres.
	genreURL = "https://api.deezer.com/radio/genres"

	musicFilePath = ""
)

// NewMusicClient returns an initialized MusicClient
// containing a list of all available genres from deezer
func NewMusicClient() *MusicClient {
	var err error
	c := &MusicClient{}
	c.Genres, err = getRadioGenres()
	if err != nil {
		panic("could not initialize music client: " + err.Error())
	}
	return c
}

func (mc *MusicClient) Play(v2jr conversation.V2JResponse) (io.ReadCloser, string, string, error) {
	isGenre := v2jr.Entities[0].Entity != "song"
	var song, trackName, artistName string
	var audioReader io.ReadCloser
	var err error

	if isGenre {
		genre := v2jr.Entities[0].Value

		trackListName, genreFound := mc.hasGenre(genre)
		if !genreFound {
			return nil, "", "", fmt.Errorf("I couldn't fine a genre by the name of %s", genre)
		}

		audioReader, trackName, artistName, err = getRandomTrack(trackListName, genre)
		if err != nil {
			return nil, "", "", err
		}
		song = trackName
	} else {
		song = strings.ReplaceAll(v2jr.Entities[0].Value, " ", "_") + ".mp3"
		f, err := ReadLocalFile(song)
		if err != nil {
			return nil, "", "", err
		}
		audioReader = io.NopCloser(f)
		trackName = song
	}

	mc.SongName = trackName
	return audioReader, trackName, artistName, nil
}

func (mc *MusicClient) PlaySpeaker(f io.ReadCloser) {
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	ctl := &beep.Ctrl{Streamer: streamer, Paused: false}
	mc.Playing = true
	mc.Ctl = ctl
	defer streamer.Close()
	done := make(chan bool)
	speaker.Play(beep.Seq(ctl, beep.Callback(func() {
		done <- true
	})))
	<-done
	mc.Playing = false
	mc.Ctl = nil
}

func (mc *MusicClient) hasGenre(genre string) (string, bool) {
	genre = strings.ToLower(genre)

	for _, e := range mc.Genres.Data {
		if strings.ToLower(e.Title) == genre || strings.Contains(strings.ToLower(e.Title), genre) {
			for _, n := range e.Radios { // look for optimal radio
				if n.Title == "Hits" || strings.ToLower(n.Title) == genre {
					return n.Tracklist, true
				}
			}
			// we couldn't find a 'pure' radio, just take the first.
			// This means the returned radio may include other genres
			return e.Radios[0].Tracklist, true
		}
	}
	return "", false
}

func (mc *MusicClient) Stop() {
	if mc.Ctl == nil {
		return
	}
	speaker.Lock()
	mc.Ctl.Paused = true
	speaker.Unlock()
	mc.Ctl = nil
	mc.Playing = false
}

func (mc *MusicClient) TogglePause() {
	if mc.Ctl == nil {
		return
	}
	speaker.Lock()
	mc.Ctl.Paused = !mc.Ctl.Paused
	mc.Playing = !mc.Ctl.Paused
	speaker.Unlock()
}

func numberInRange(lower, upper int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn((upper - lower + 1) + lower)
}

func getRadioGenres() (RadioGenreResponse, error) {
	genres := RadioGenreResponse{}
	resp, err := http.DefaultClient.Get(genreURL)
	if err != nil {
		return genres, fmt.Errorf("could not get music genres: %s", err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&genres)
	if err != nil {
		return genres, fmt.Errorf("could not decode genre response: %s", err.Error())
	}
	return genres, nil
}

func getTracks(trackList string) (RadioTrackListResponse, error) {
	tracks := RadioTrackListResponse{}
	resp, err := http.DefaultClient.Get(trackList)
	if err != nil {
		return tracks, fmt.Errorf("could not get tracks: " + err.Error())
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tracks)
	if err != nil {
		return tracks, fmt.Errorf("could not handle tracks response: " + err.Error())
	}
	return tracks, nil
}

func getRandomTrack(trackListName, genre string) (io.ReadCloser, string, string, error) {
	trackList, err := getTracks(trackListName)
	if err != nil {
		fmt.Println(err)
		return nil, "", "", fmt.Errorf("sorry, I had trouble getting a list of tracks for the %s genre", genre)
	}

	trackNum := numberInRange(0, len(trackList.Data)-1)

	track := trackList.Data[trackNum]

	audio, err := DownloadPreview(track.Preview)
	if err != nil {
		fmt.Println(err)
		return nil, "", "", fmt.Errorf("sorry, I had trouble downloading the requested music")
	}
	return io.NopCloser(bytes.NewReader(audio)), track.Title, track.Artist.Name, nil
}

func DownloadPreview(url string) ([]byte, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return by, nil
}

func ReadLocalFile(song string) (*os.File, error) {
	_, err := os.Stat(musicFilePath + song)
	if err != nil {
		fmt.Printf("file does not exist: %s\n", err.Error())
		return nil, fmt.Errorf("sorry, I couldn't find that song")
	}

	f, err := os.Open(musicFilePath + song)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("sorry, I couldn't open that song")
	}
	return f, nil
}
