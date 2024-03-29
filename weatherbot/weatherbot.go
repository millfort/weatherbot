// Package weatherbot implements a simple library for creating telegram bot that shows the current weather.
package weatherbot

import (
	"fmt"

	owm "github.com/briandowns/openweathermap"
	"github.com/yanzay/tbot/v2"
)

type weatherBot struct {
	*tbot.Server
	bot          *tbot.Server
	client       *tbot.Client
	weatherAPI   *owm.CurrentWeatherData
	weatherIcons map[string]string
}

// WeatherBot is the interface that provides a set of functions for using weatherbot
type WeatherBot interface {
	StartHandler(*tbot.Message)
	WeatherHandler(*tbot.Message)
	Start() error
	HandleMessage(string, func(*tbot.Message))
}

// New attempts to creates a new WeatherBot
func New(tgKey, owmKey string) (WeatherBot, error) {
	wAPI, err := owm.NewCurrent("C", "en", owmKey)
	if err != nil {
		return nil, err
	}
	server := tbot.New(tgKey)
	client := server.Client()
	weatherBot := weatherBot{
		Server:       server,
		client:       client,
		weatherAPI:   wAPI,
		weatherIcons: getWeatherIcons(),
	}
	return &weatherBot, nil
}

// StartHandler is a handler functions for initiate conversation with bot
func (wb *weatherBot) StartHandler(msg *tbot.Message) {
	wb.client.SendMessage(msg.Chat.ID, "Напишите мне название вашего города и я скажу вам температуру")
}

// WeatherHandler is a handler function for for requesting current weather
func (wb *weatherBot) WeatherHandler(msg *tbot.Message) {
	err := wb.weatherAPI.CurrentByName(msg.Text)
	if err != nil {
		wb.client.SendMessage(msg.Chat.ID, "Не могу найти ваш город")
		return
	}
	weather := wb.getWeatherIcon()
	wb.client.SendMessage(msg.Chat.ID, fmt.Sprintf("В городе %s %.2f С° %s", msg.Text, wb.weatherAPI.Main.Temp, weather))
}

func getWeatherIcons() map[string]string {
	return map[string]string{
		"Thunderstorm": "☁☔⚡",
		"Drizzle":      "☔",
		"Rain":         "☁☔",
		"Snow":         "❄",
		"Clear":        "☀",
		"Clouds":       "⛅",
	}
}

func (wb *weatherBot) getWeatherIcon() string {
	if len(wb.weatherAPI.Weather) != 0 {
		if icon, exist := wb.weatherIcons[wb.weatherAPI.Weather[0].Main]; exist {
			return icon
		}
	}
	return ""
}
