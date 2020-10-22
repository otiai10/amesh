package commands

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/otiai10/amesh/bot/slack"
	"github.com/otiai10/goapis/openweathermap"
)

// ForecastCommand ...
type ForecastCommand struct{}

// ForecastCommandError ...
type ForecastCommandError error

// ErrorForecastNoEntry ...
var ErrorForecastNoEntry ForecastCommandError = errors.New("十分な天気予報情報が取得できませんでした")

// Match ...
func (cmd ForecastCommand) Match(payload *slack.Payload) bool {
	if len(payload.Ext.Words) == 0 {
		return false
	}
	return payload.Ext.Words[0] == "予報" || payload.Ext.Words[0] == "forecast"
}

// Handle ...
func (cmd ForecastCommand) Handle(ctx context.Context, payload *slack.Payload) slack.Message {
	client := openweathermap.New(os.Getenv("OPENWEATHERMAP_API_KEY"))
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return wrapError(payload, err)
	}
	city := "Tokyo"
	res, err := client.ByCityName(city, nil)
	if err != nil {
		return wrapError(payload, err)
	}
	if len(res.Forecasts) == 0 || len(res.Forecasts[0].Weather) == 0 {
		return wrapError(payload, ErrorForecastNoEntry)
	}

	message := slack.Message{
		Channel: payload.Event.Channel,
		Text:    res.City.Name,
	}

	// {{{ 日付で分けて、行をつくっていく
	var blockdate int
	placeholder := cmd.getPlaceholderEmoji()
	for i, forecast := range res.Forecasts {
		w := forecast.Weather[0]
		t := time.Unix(forecast.Timestamp, 0).In(loc)
		_, month, date := t.Date()
		// 新しい日付であればBlockを初期化
		if date != blockdate {
			message.Text += "\n"
			message.Text += fmt.Sprintf("%d/%d %s ", month, date, cmd.getJapaneseWeekday(t.Weekday()))
			if t.Hour() != 0 {
				for h := 0; h < t.Hour(); h += 3 {
					message.Text += placeholder
				}
			}
			blockdate = date
		}
		emoji := cmd.convertOpenWeatherMapIconToSlackEmoji(w.Icon)
		message.Text += emoji
		if i == len(res.Forecasts)-1 && t.Hour() != 21 {
			for h := t.Hour() + 3; h < 24; h += 3 {
				message.Text += placeholder
			}
		}
	}
	// }}}

	return message
}

// Help ...
func (cmd ForecastCommand) Help(payload *slack.Payload) slack.Message {
	return slack.Message{
		Channel: payload.Event.Channel,
		Text:    "天気予報コマンド",
	}
}

func (cmd ForecastCommand) convertOpenWeatherMapIconToSlackEmoji(icon string) string {
	// https://openweathermap.org/weather-conditions
	dictionary := map[string]string{
		"01": ":sunny:",
		"02": ":mostly_sunny:",
		"03": ":partly_sunny:",
		"04": ":cloud:",
		"09": ":rain_cloud:",
		"10": ":partly_sunny_rain:",
		"11": ":thunder_cloud_and_rain:",
		"13": ":snowflake:",
		"50": ":fog:",
	}
	if emoji, ok := dictionary[icon[:2]]; ok {
		return emoji
	}
	return cmd.getPlaceholderEmoji()
}

func (cmd ForecastCommand) getPlaceholderEmoji() string {
	candidates := []string{
		":marijuana:",
		":shrimp:",
		":pig:",
		":slack:",
		":sunglasses:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
		":white_small_square:",
	}
	return candidates[rand.Intn(len(candidates))]
}

func (cmd ForecastCommand) getJapaneseWeekday(day time.Weekday) string {
	return map[time.Weekday]string{
		time.Sunday:    "日",
		time.Monday:    "月",
		time.Tuesday:   "火",
		time.Wednesday: "水",
		time.Thursday:  "木",
		time.Friday:    "金",
		time.Saturday:  "土",
	}[day]
}
