package conversation

// add more intents as sentences.ini grows
const (
	CurrentWeatherIntent  = "Weather"
	DailyWeatherForecast  = "DailyWeatherForecast"
	HourlyWeatherForecast = "HourlyWeatherForecast"
	TimeIntent            = "GetTime"
	HelpIntent            = "Help"
	RancherIntent         = "RancherInfo"
	TemperatureIntent     = "GetTemperature"
	PlayMusicIntent       = "PlayMusic"
	StopMusicIntent       = "StopMusic"
)

type V2JResponse struct {
	Text              string   `json:"text"`
	Likelihood        float64  `json:"likelihood"`
	TranscribeSeconds float64  `json:"transcribe_seconds"`
	WavSeconds        float64  `json:"wav_seconds"`
	Tokens            []string `json:"tokens"`
	Timeout           bool     `json:"timeout"`
	Intent            struct {
		Name       string  `json:"name"`
		Confidence float64 `json:"confidence"`
	} `json:"intent"`
	Entities []struct {
		Entity    string   `json:"entity"`
		Value     string   `json:"value"`
		RawValue  string   `json:"raw_value"`
		Source    string   `json:"source"`
		Start     float64  `json:"start"`
		RawStart  float64  `json:"raw_start"`
		End       float64  `json:"end"`
		RawEnd    float64  `json:"raw_end"`
		Tokens    []string `json:"tokens"`
		RawTokens []string `json:"raw_tokens"`
	} `json:"entities"`
	RawText          string      `json:"raw_text"`
	RecognizeSeconds float64     `json:"recognize_seconds"`
	RawTokens        []string    `json:"raw_tokens"`
	SpeechConfidence float64     `json:"-"` // doesn't seem to work
	WavName          interface{} `json:"wav_name"`
	Slots            struct {
		State string `json:"state"`
		Name  string `json:"name"`
	} `json:"slots"`
}
