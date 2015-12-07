package amesh

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"time"
)

// EventType ...
type EventType string

const (
	// Start ...
	Start EventType = "start"
	// Error ...
	Error EventType = "error"
	// Update ...
	Update EventType = "update"
	// Rain ...
	Rain EventType = "rain"
)

// Observer ...
type Observer struct {
	handlers          map[EventType]EventHandleFunc
	IterationDuration time.Duration
	IsRaining         func(ev Event) bool
	Notifier          Notifier
}

// Event ...
type Event struct {
	Error     error
	Img       image.Image
	Timestamp time.Time
}

// EventHandleFunc ...
type EventHandleFunc func(Event) error

// NewObserver ...
func NewObserver() *Observer {
	return &Observer{
		handlers: map[EventType]EventHandleFunc{
			Update: func(event Event) error {
				log.Println("[UPDATE]", event.Timestamp)
				return nil
			},
			Start: func(event Event) error {
				log.Println("[START]", event.Timestamp)
				return nil
			},
			Error: func(event Event) error {
				panic(event)
			},
			Rain: DefaultOnRainHandleFunc,
		},
		IterationDuration: DefaultIterationDuration,
		// Set custom rain judgement func here.
		IsRaining: DefaultIsRainingFunc,
	}
}

// On ...
func (observer *Observer) On(eventtype EventType, fun EventHandleFunc) *Observer {
	switch eventtype {
	default:
		observer.handlers[eventtype] = fun
	}
	return observer
}

// Start ...
func (observer *Observer) Start() {

	onerror := make(chan Event)

	go func() {
		for {
			select {
			case ev := <-onerror:
				observer.handlers[Error](ev)
			}
		}
	}()

	observer.handlers[Start](Event{Timestamp: time.Now()})
	observer.loop(onerror)
}

// Restart is just an alias for Start.
func (observer *Observer) Restart() {
	observer.Start()
}

func (observer *Observer) loop(onerror chan Event) {

	ticker := time.Tick(observer.IterationDuration)

	for {
		<-ticker
		err := observer.Run()
		if err != nil {
			onerror <- Event{Error: fmt.Errorf("%v", err)}
			break
		}
	}
}

// Run ...
func (observer *Observer) Run() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	entry := GetEntry()
	res, err := http.Get(entry.Mesh)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return err
	}

	event := Event{
		Img:       img,
		Timestamp: time.Now(),
	}

	if _, ok := observer.handlers[Rain]; ok && observer.IsRaining(event) {
		err = observer.handlers[Rain](event)
	} else {
		err = observer.handlers[Update](event)
	}

	return err
}

// SetNotifier ...
func (observer *Observer) SetNotifier(notifier Notifier) *Observer {
	observer.Notifier = notifier
	return observer
}
