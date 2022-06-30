package agent

import (
	"fmt"
	"github.com/spf13/viper"
	"harrisonwaffel/assistant/pkg/conversation"
	"harrisonwaffel/assistant/pkg/utils"
)

func handleIntent(ctx *agentContext, resp conversation.V2JResponse) {
	if CurrentlyThinking {
		return
	}
	CurrentlyThinking = true
	switch resp.Intent.Name {
	case conversation.DailyWeatherForecast:
		ctx.speak(ctx.WeatherCache.ForecastDailyWeather(resp))

	case conversation.HourlyWeatherForecast:
		ctx.speak(ctx.WeatherCache.ForecastHourlyWeather(resp)...)

	case conversation.CurrentWeatherIntent:
		ctx.speak(ctx.WeatherCache.GeneralWeather(resp))

	case conversation.TemperatureIntent:
		ctx.speak(ctx.WeatherCache.Temperature(resp))

	case conversation.TimeIntent:
		ctx.speak(utils.GetTimeAsSentence())

	case conversation.RancherIntent:
		if viper.GetBool("apps.rancher.enabled") {
			// do rancher stuff
		}
	case conversation.PlayMusicIntent:
		audioReader, trackName, artistName, errMessage := ctx.MusicClient.Play(resp)
		if errMessage != nil {
			ctx.speak(errMessage.Error())
		}

		if artistName != "" {
			ctx.speak(fmt.Sprintf("OK, playing %s by %s", trackName, artistName))
		} else {
			ctx.speak(fmt.Sprintf("OK, playing %s", trackName))
		}

		ctx.MusicClient.PlaySpeaker(audioReader)

	case conversation.StopMusicIntent:
		ctx.MusicClient.Stop()

	case conversation.HelpIntent:
		ctx.speak("I am Geeko, a very basic open-source voice assistant. You can ask me for the weather, the time, or to play music.")
	}

	CurrentlyThinking = false
}

func (ctx *agentContext) speak(sentence ...string) {
	if ctx.PrettyPrint {
		go func() {
			for _, e := range sentence {
				WordChan <- e
			}
		}()
	}
	ctx.Voice2JsonCommand.speak(sentence...)
}

func (v Voice2JsonCommand) speak(sentence ...string) {
	utils.Speak(sentence...)
}
