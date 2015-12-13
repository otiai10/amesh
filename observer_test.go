package amesh

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "image/gif"
)

func TestNewObserver(t *testing.T) {

	observer := NewObserver()

	count := 1

	observer.IterationDuration = 2 * time.Second
	observer.IsRaining = func(ev Event) bool {
		return (count%3 == 0)
	}

	observer.On(Rain, func(ev Event) error {
		count++
		log.Println("[test][RAIN]", ev)
		return nil
	})
	observer.On(Update, func(ev Event) error {
		if count > 5 {
			return fmt.Errorf("なんかエラー")
		}
		count++
		log.Println("[test][UPDATE]", ev)
		return nil
	})
	observer.On(Error, func(ev Event) error {
		log.Println("[test][ERROR]", "エラーが起きた")
		return nil
	})

	observer.Start()
}
