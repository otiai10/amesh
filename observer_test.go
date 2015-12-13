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
		return nil
	})
	observer.On(Update, func(ev Event) error {
		if count > 5 {
			return fmt.Errorf("なんかエラー")
		}
		count++
		return nil
	})
	observer.On(Error, func(ev Event) error {
		log.Println("[ERROR]", ev.Timestamp)
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
		return nil
	})

	go observer.Start()
	time.Sleep(4 * time.Second)
	observer.Stop()

	Expect(t, count).ToBe(4)
}

func TestObserver_NotificationInterval(t *testing.T) {
	observer := NewObserver(1 * time.Second)
	observer.IsRaining = func(ev Event) bool {
		return true
	}

	dummy := &DummyNotifier{}
	observer.Notifier = dummy
	observer.NotificationInterval = 3 * time.Second

	observer.On(Rain, func(ev Event) error {
		msg := fmt.Sprintf("雨")
		if observer.LastRain.IsZero() && observer.Notifier != nil {
			if err := observer.Notifier.Notify(msg); err != nil {
				return err
			}
		}
		if observer.LastRain.IsZero() {
			observer.LastRain = ev.Timestamp
		}
		if ev.Timestamp.After(observer.LastRain.Add(observer.NotificationInterval)) {
			observer.LastRain = time.Time{} // reset to notify again
		}
		return nil
	})

	go observer.Start()
	time.Sleep(8 * time.Second)
	observer.Stop()

	Because(t, "Notification should be throttled by its interval", func(t *testing.T) {
		Expect(t, dummy.count).ToBe(2)
	})
}

type DummyNotifier struct {
	last  string
	count int
}

func (n *DummyNotifier) Notify(msg string) error {
	n.last = msg
	n.count++
	return nil
}
