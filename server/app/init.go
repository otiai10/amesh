package app

import (
	"net/http"

	"github.com/otiai10/amesh/server/services"
	"github.com/otiai10/marmoset"
)

func init() {
	router := marmoset.NewRouter()
	router.POST("/webhook", services.Handler().ServeHTTP)
	http.Handle("/", router)
}
