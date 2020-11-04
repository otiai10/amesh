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

	if payload.Ext.Words.Flag("-h") {
		return cmd.Help(payload)
	}

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
		Text:    res.City.Name + "\n",
	}

	placeholder := cmd.getPlaceholderEmoji()
	for _, group := range res.GroupByDate(loc) {
		message.Text += ForecastRowBuilder{
			Placeholder: placeholder,
			IncludTemp:  payload.Ext.Words.Flag("-t"),
		}.build(group) + "\n"
	}

	return message
}

// Help ...
func (cmd ForecastCommand) Help(payload *slack.Payload) slack.Message {
	return slack.Message{
		Channel: payload.Event.Channel,
		Text:    "天気予報コマンド\n```@amesh [予報|forecast] [-t|-h]```",
	}
}

func (cmd ForecastCommand) getPlaceholderEmoji() string {
	rand.Seed(time.Now().Unix())
	if rand.Intn(50) > 1 {
		return ":white_small_square:"
	}
	candidates := []string{
		":marijuana:",
		":shrimp:",
		":pig:",
		":slack:",
		":sunglasses:",
	}
	return candidates[rand.Intn(len(candidates))]
}

// ForecastRowBuilder ...
type ForecastRowBuilder struct {
	// 日付
	Month   time.Month
	Date    int
	Weekday time.Weekday

	// 時間ごとの天気
	// Weather []openweathermap.Weather
	Head        string
	Body        string
	Placeholder string

	// 気温
	IncludTemp bool
	MaxCelsius float32
	MinCelsius float32
}

func (row ForecastRowBuilder) build(forecast []openweathermap.Forecast) string {
	first := forecast[0]
	last := forecast[len(forecast)-1]
	_, month, date := first.LocalTime.Date()
	row.Month = month
	row.Date = date
	row.Weekday = first.LocalTime.Weekday()
	row.MinCelsius = first.Main.TempMin
	row.MaxCelsius = first.Main.TempMax
	row.Head = fmt.Sprintf("%02d/%02d %s ", row.Month, row.Date, row.getJapaneseWeekday())
	for h := 0; h < first.LocalTime.Hour(); h += 3 {
		row.Body += row.Placeholder
	}
	for _, f := range forecast {
		row.Body += row.convertOpenWeatherMapIconToSlackEmoji(f.Weather[0].Icon)
		if f.Main.TempMin < row.MinCelsius {
			row.MinCelsius = f.Main.TempMin
		}
		if f.Main.TempMax > row.MaxCelsius {
			row.MaxCelsius = f.Main.TempMax
		}
	}
	for h := 21; h > last.LocalTime.Hour(); h -= 3 {
		row.Body += row.Placeholder
	}
	if row.IncludTemp {
		return fmt.Sprintf("%s %s %v/%v", row.Head, row.Body, row.MinCelsius, row.MaxCelsius)
	}
	return fmt.Sprintf("%s %s", row.Head, row.Body)
}

func (row ForecastRowBuilder) convertOpenWeatherMapIconToSlackEmoji(icon string) string {
	// https://openweathermap.org/weather-conditions
	dictionary := map[string]string{
		"01": ":sunny:",
		"02": ":mostly_sunny:",
		"03": ":partly_sunny:",
		"04": ":cloud:",
		"09": ":rain_cloud:",
		"10": ":umbrella:",
		"11": ":umbrella_with_rain_drops:",
		"13": ":snowflake:",
		"50": ":fog:",
	}
	if emoji, ok := dictionary[icon[:2]]; ok {
		return emoji
	}
	return row.Placeholder
}

func (row ForecastRowBuilder) getJapaneseWeekday() string {
	return map[time.Weekday]string{
		time.Sunday:    "*日*",
		time.Monday:    "月",
		time.Tuesday:   "火",
		time.Wednesday: "水",
		time.Thursday:  "木",
		time.Friday:    "金",
		time.Saturday:  "*土*",
	}[row.Weekday]
}
