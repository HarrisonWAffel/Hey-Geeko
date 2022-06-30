package weather

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"harrisonwaffel/assistant/pkg/conversation"
	"harrisonwaffel/assistant/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var WeatherCache Cache

func init() {
	WeatherCache.responses = make(map[string]CacheItem)
}

const (
	URL = "https://api.weatherapi.com/v1/forecast.json?"
)

func buildURL(location string) string {
	return fmt.Sprintf("%skey=%s&q=%s&days=1&api=no&alerts=no", URL, viper.GetString("apps.weather.key"), location)
}

func getLocation(v2jr conversation.V2JResponse) string {
	location := ""
	if len(v2jr.Entities) == 0 {
		// assume default location
		location = viper.GetString("apps.weather.defaultLocation")
	} else {
		location = strings.ReplaceAll(v2jr.Entities[0].Value, " ", "_")
	}
	return location
}

func (c *Cache) getWeather(v2j conversation.V2JResponse) (Response, string) {

	if viper.GetString("apps.weather.key") == "" {
		return Response{}, "Please setup a weatherapi.com API key"
	}

	location := getLocation(v2j)
	resp, responseCached := c.responses[location]
	if !responseCached {
		c.fetchWeather(location)
		return c.responses[location].Response, ""
	}

	if resp.gotten.Add(10 * time.Minute).Before(time.Now()) {
		c.fetchWeather(location)
	}

	return c.responses[location].Response, ""
}

func (c *Cache) fetchWeather(location string) {
	resp, err := http.DefaultClient.Get(buildURL(location))
	if err != nil {
		panic("could not get weather: " + err.Error())
	}

	weatherResp := Response{}
	err = json.NewDecoder(resp.Body).Decode(&weatherResp)
	if err != nil {
		panic("could not get weather: " + err.Error())
	}

	item := CacheItem{
		gotten:   time.Now(),
		Response: weatherResp,
	}
	c.responses[location] = item
}

// GeneralWeather gets the current, or cached, weather API response
// for the given location and reads back the current conditions and temperatures
func (c *Cache) GeneralWeather(v2jr conversation.V2JResponse) string {
	resp, errResp := c.getWeather(v2jr)
	if errResp != "" {
		return errResp
	}
	location := getLocation(v2jr)
	return fmt.Sprintf("In %s it is %s. %s.", location, resp.Current.Condition.Text, c.Temperature(v2jr))
}

// Temperature gets the current, or cached, weather API response
// for the given location and reads back the temperature
func (c *Cache) Temperature(v2jr conversation.V2JResponse) string {
	resp, errResp := c.getWeather(v2jr)
	if errResp != "" {
		return errResp
	}
	return fmt.Sprintf("It is currently %.1f degrees, and it feels like %.1f degrees", resp.Current.TempF, resp.Current.FeelslikeF)
}

// ForecastDailyWeather gets the current, or cached, weather API response
// and reads the temperature for the next three hours
func (c *Cache) ForecastDailyWeather(v2jr conversation.V2JResponse) string {
	weather, errResp := c.getWeather(v2jr)
	if errResp != "" {
		return errResp
	}

	// get the weather for the next three hours
	if len(weather.Forecast.Forecastday) == 0 {
		return "Could not get daily weather forecast"
	}

	forecast := weather.Forecast.Forecastday[0]
	return fmt.Sprintf("it is %s. With a high of %.2f and a low of %.2f", forecast.Day.Condition, forecast.Day.MaxtempF, forecast.Day.MintempF)
}

func (c *Cache) ForecastHourlyWeather(v2jr conversation.V2JResponse) []string {
	weather, errResp := c.getWeather(v2jr)
	if errResp != "" {
		return []string{errResp}
	}
	if len(weather.Forecast.Forecastday) == 0 {
		return []string{"Could not get daily weather forecast"}
	}

	var sentences []string
	format := "15:04"
	for _, hour := range weather.Forecast.Forecastday[0].Hour {
		thisHour := strings.Split(hour.Time, " ")[1]
		t, err := time.Parse(format, thisHour)
		if err != nil {
			continue
		}

		// we don't need past reports
		if t.Hour() < time.Now().In(utils.GetLoc()).Hour() {
			continue
		}

		thi, err := strconv.Atoi(strings.Split(thisHour, ":")[0])
		if err != nil {
			fmt.Println(err)
			continue
		}

		amPm := "AM"
		if thi >= 12 {
			amPm = "PM"
			if thi != 12 {
				thi -= 12
			}
		}

		thisHour = fmt.Sprintf("%d %s", thi, amPm)

		forecast := fmt.Sprintf("At %s in %s it will be %s and %.1f degrees. ", thisHour, getLocation(v2jr), hour.Condition.Text, hour.TempF)
		if hour.ChanceOfSnow != 0 {
			forecast = fmt.Sprintf("%s and There is a %d chance of snow.", forecast, hour.ChanceOfSnow)
		}
		if hour.ChanceOfRain != 0 {
			forecast = fmt.Sprintf("%s and There is a %d chance of rain.", forecast, hour.ChanceOfRain)
		}
		sentences = append(sentences, forecast)
		if len(sentences) == 3 {
			break
		}
	}

	return sentences
}
