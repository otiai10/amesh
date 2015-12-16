package amesh

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"time"

	_ "image/gif"
)

// EventType ...
type EventType string

const (
	// Start ...
	Start EventType = "start"
	// Stop ...
	Stop EventType = "stop"
	// Error ...
	Error EventType = "error"
	// Update ...
	Update EventType = "update"
	// Rain ...
	Rain EventType = "rain"
)

// Observer ...
type Observer struct {
	handlers             map[EventType]EventHandleFunc
	IterationDuration    time.Duration
	IsRaining            func(ev Event) bool
	Notifier             Notifier
	NotificationInterval time.Duration
	LastRain             time.Time
	onerror              chan Event
	stopper              chan Event
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
func NewObserver(duration ...time.Duration) *Observer {
	duration = append(duration, DefaultIterationDuration)
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
			Stop: func(event Event) error {
				log.Println("[STOP]", event.Timestamp)
				return nil
			},
			Error: func(event Event) error {
				// panic(event.Error)
				log.Println("[ERROR]", event)
				return nil
			},
			Rain: DefaultOnRainHandleFunc,
		},
		IterationDuration: duration[0],
		// Set custom rain judgement func here.
		IsRaining: DefaultIsRainingFunc,
		onerror:   make(chan Event),
		stopper:   make(chan Event),
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

// Start never ends unless error or stop called
func (observer *Observer) Start() {

	observer.handlers[Start](Event{Timestamp: time.Now()})
	go observer.loop()

	ev := <-observer.onerror
	observer.handlers[Error](ev)
}

// Stop ...
func (observer *Observer) Stop() {
	observer.stopper <- Event{Timestamp: time.Now()}
}

// Restart is just an alias for Start.
func (observer *Observer) Restart() {
	observer.Start()
}

func (observer *Observer) loop() {

	ticker := time.Tick(observer.IterationDuration)

	for {
		select {
		case <-observer.stopper:
			observer.onerror <- Event{Error: fmt.Errorf("%v", "interrupted")}
			return
		case <-ticker:
			err := observer.Run()
			if err != nil {
				observer.onerror <- Event{Error: fmt.Errorf("%v", err)}
				return
			}
		}
	}
}

// Run ...
func (observer *Observer) Run() error {

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
