package main

import (
	"log"
	"net/http"
	"os"

	"github.com/otiai10/amesh/bot/slack"
	"github.com/otiai10/marmoset"
)

func init() {
	router := marmoset.NewRouter()
	router.POST("/slack/webhook", slack.HandleWebhook)
	router.GET("/slack", slack.HandleIndex)
	http.Handle("/", router)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
