package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/otiai10/amesh/bot/slack"
	"gopkg.in/yaml.v2"
)

func main() {
	devLoadEnv("./bot/app-secrets.yaml")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	u, _ := url.Parse(fmt.Sprintf("http://localhost:%s", port))
	u.Path = "/slack/webhook"
	body := bytes.NewBuffer(nil)

	text := strings.Join(append([]string{"@amesh"}, os.Args[1:]...), " ")
	channel := "otiai10-dev"
	json.NewEncoder(body).Encode(slack.Payload{
		Token: os.Getenv("SLACK_VERIFICATION_TOKEN"),
		Event: slack.Event{
			Text:    text,
			Channel: channel,
		},
	})
	_, err := http.Post(u.String(), "application/json", body)
	fmt.Println(err)
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
