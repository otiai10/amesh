package amesh

import (
	"fmt"
	"log"
	"testing"
	"time"

	_ "image/gif"

	. "github.com/otiai10/mint"
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

	Expect(t, count).ToBe(7)
}

func TestObserver_Stop(t *testing.T) {

	observer := NewObserver(1 * time.Second)
	observer.IsRaining = func(ev Event) bool {
		return false
	}

	count := 1
	observer.On(Update, func(ev Event) error {
		count++
		log.Println(count)
		return nil
	})

	go observer.Start()
	time.Sleep(4 * time.Second)
	observer.Stop()

	Expect(t, count).ToBe(4)
}
