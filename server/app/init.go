package app

import (
	"net/http"
	"os"

	"github.com/otiai10/amesh/server/plugins"
	"github.com/otiai10/amesh/server/services"
	"github.com/otiai10/marmoset"
)

func init() {
	router := marmoset.NewRouter()

	p := []plugins.Plugin{
		plugins.Image{GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"), GoogleCustomSearchEngineID: os.Getenv("GOOGLE_CUSTOM_SEARCH_ENGINE_ID")},
		plugins.Gif{GoogleAPIKey: os.Getenv("GOOGLE_API_KEY"), GoogleCustomSearchEngineID: os.Getenv("GOOGLE_CUSTOM_SEARCH_ENGINE_ID")},
		plugins.Meshi{YelpAPIKey: os.Getenv("YELP_API_KEY")},
		plugins.WeatherForecast{APIKey: os.Getenv("OPENWEATHERMAP_API_KEY")},
	}

	s := services.Handler(p...)
	router.POST(s.WebhookURL(), s.HandleWebhook)
	router.POST(s.QueueURL(), s.HandleQueue)

	http.Handle("/", router)
}
