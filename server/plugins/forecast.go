package plugins

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/otiai10/amesh/server/middlewares"
	"github.com/otiai10/openweathermap"
)

// WeatherForecast ...
// OpenWeatherMap API を使って天気予報を返す
type WeatherForecast struct {
	APIKey string
}

// Method ...
func (wf WeatherForecast) Method() string {
	return "forecast"
}

// Match ...
func (wf WeatherForecast) Match(ctx context.Context, texts []string) bool {
	if len(texts) == 0 {
		return false
	}
	return texts[0] == "forecast"
}

// TaskValues ...
func (wf WeatherForecast) TaskValues(ctx context.Context, texts []string) url.Values {
	return url.Values{"query": {strings.Join(texts[1:], ",")}}
}

// Exec ...
// TODO: たぶんstringじゃないほうがいいんだよねえ
func (wf WeatherForecast) Exec(ctx context.Context, r *http.Request) (string, error) {

	client := openweathermap.New(wf.APIKey)
	client.HTTPClient = middlewares.HTTPClient(ctx)

	q := r.FormValue("query")

	res, err := client.ByCityName(q, nil)
	if err != nil {
		return "", err
	}

	// TODO: Execがstringを返すんじゃなくてもっとおもしろいものを返すようになるべき
	lines := []string{fmt.Sprintf("City: %s", res.City.Name)}
	for _, forecast := range res.Forecasts {
		lines = append(lines, fmt.Sprintf("%s\t%s", forecast.Datetime, forecast.Weather[0].Description))
	}

	return strings.Join(lines, "\n"), nil
}
