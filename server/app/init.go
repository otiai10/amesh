package app

import (
	"net/http"

	"github.com/otiai10/amesh/server/services"
	"github.com/otiai10/marmoset"
)

func init() {
	router := marmoset.NewRouter()

	s := services.Handler()
	router.POST(s.WebhookURL(), s.HandleWebhook)
	router.POST(s.QueueURL(), s.HandleQueue)

	http.Handle("/", router)
}
