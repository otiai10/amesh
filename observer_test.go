package amesh

import (
	"log"
	"testing"

	_ "image/gif"
)

func TestNewObserver(t *testing.T) {

	observer := NewObserver()

	observer.IsRaining = func(ev Event) bool {
		// return true
		return true
	}
	observer.On(Rain, func(ev Event) {
		log.Println(ev)
	})
	observer.On(Update, func(ev Event) {
		log.Println("ほげ", ev)
	})
	observer.On(Error, func(ev Event) {
		log.Println("エラーが起きた")
	})

	observer.Start()
}
