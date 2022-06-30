package weather

import "time"

type Cache struct {
	responses map[string]CacheItem
}

type CacheItem struct {
	Response
	gotten time.Time
}

type WeatherRequest struct {
	City string
}

type Response struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempF       float64 `json:"temp_f"`
		IsDay       int     `json:"is_day"`
		Condition   struct {
			Text string `json:"text"`
		} `json:"condition"`
		Cloud      int     `json:"cloud"`
		FeelslikeF float64 `json:"feelslike_f"`
		Uv         float64 `json:"uv"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Day  struct {
				MaxtempF  float64 `json:"maxtemp_f"`
				MintempF  float64 `json:"mintemp_f"`
				AvgtempF  float64 `json:"avgtemp_f"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"day"`
			Astro struct {
				Sunrise string `json:"sunrise"`
				Sunset  string `json:"sunset"`
			} `json:"astro"`
			Hour []struct {
				Time      string  `json:"time"`
				TempF     float64 `json:"temp_f"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				Humidity     int     `json:"humidity"`
				FeelslikeF   float64 `json:"feelslike_f"`
				WindchillF   float64 `json:"windchill_f"`
				HeatindexF   float64 `json:"heatindex_f"`
				DewpointF    float64 `json:"dewpoint_f"`
				WillItRain   int     `json:"will_it_rain"`
				ChanceOfRain int     `json:"chance_of_rain"`
				WillItSnow   int     `json:"will_it_snow"`
				ChanceOfSnow int     `json:"chance_of_snow"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}
