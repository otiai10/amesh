package app

import (
	"net/http"

	"github.com/otiai10/amesh/server/controllers"
	"github.com/otiai10/marmoset"
)

func init() {
	r := marmoset.NewRouter()
	r.POST("/webhook", controllers.Webhook)
	http.Handle("/", r)
}
