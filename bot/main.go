package main

import (
	"log"
	"net/http"
	"os"

	"github.com/otiai10/amesh/bot/slack"
	"github.com/otiai10/marmoset"
	"gopkg.in/yaml.v2"
)

func init() {
	router := marmoset.NewRouter()
	router.POST("/slack/webhook", slack.HandleWebhook)
	router.GET("/slack", slack.HandleIndex)
	http.Handle("/", router)
}

func main() {

	if os.Getenv("GAE_APPLICATION") == "" {
		devLoadEnv("./bot/app-secrets.yaml")
	}

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

func devLoadEnv(fname string) {
	// AppConfig ...
	type AppConfig struct {
		EnvVariables map[string]string `yaml:"env_variables"`
	}
	log.Printf("Loading env variables from %s", fname)
	f, err := os.Open(fname)
	if err != nil {
		log.Fatalln(err)
	}
	appconfig := &AppConfig{}
	if err := yaml.NewDecoder(f).Decode(appconfig); err != nil {
		log.Fatalln(err)
	}
	for name, value := range appconfig.EnvVariables {
		os.Setenv(name, value)
	}
}
