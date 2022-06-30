package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func GetTimeAsSentence() string {
	locString := viper.GetString("apps.time.location")
	if locString == "" {
		locString = "America/New_York"
	}
	loc, err := time.LoadLocation(locString)
	if err != nil {
		panic("Could not load location: " + err.Error())
	}
	return fmt.Sprintf("It is %s", time.Now().In(loc).Format(time.Kitchen))
}

func GetLoc() *time.Location {
	locString := viper.GetString("apps.time.location")
	if locString == "" {
		locString = "America/New_York"
	}
	loc, err := time.LoadLocation(locString)
	if err != nil {
		panic(err)
	}
	return loc
}
