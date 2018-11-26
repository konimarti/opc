package opc

import (
	// "log"
	// "errors"
	"io"
	"time"
)

//Trigger interface to implement when to notify observers
type Trigger interface {
	Fire(interface{}) bool
	Update(interface{})
}

//OnChange implements the Trigger interface.
//Notifies observers when value changes.
type OnChange struct {
	value interface{}
}

//Fire checks if observers should be notified
func (t *OnChange) Fire(newValue interface{}) bool { return t.value != newValue }

//Update updates internally stored value
func (t *OnChange) Update(newValue interface{}) { t.value = newValue }

//OnValue implements the Trigger interface.
//Notifies observers when a new value matched to stored value.
type OnValue struct {
	value interface{}
}

//Fire checks if observers should be notified
func (t *OnValue) Fire(newValue interface{}) bool { return t.value == newValue }

//Update updates internally stored value
func (t *OnValue) Update(newValue interface{}) {}

//Observer observes and notified based on the triggers.
//Note the observer should be replaced with the observer package
//in the near future (github.com/konimarti/observer).
type Observer struct {
	trigger   Trigger
	conn      Connection
	observers []chan bool
	closing   []*control
}

//Close gracefully shutdowns observer
func (o *Observer) Close() {
	for _, control := range o.closing {
		control.Close()
	}
}

//Notify notified registered observers
func (o *Observer) Notify() {
	for _, observer := range o.observers {
		select {
		case <-observer:
		default:
		}
		observer <- true
	}
}

//Channel provides channel for observation
func (o *Observer) Channel() chan bool {
	observer := make(chan bool, 1)
	o.observers = append(o.observers, observer)
	return observer
}

//Observe starts the observation
func (o *Observer) Observe(tag string, refresh time.Duration) io.Closer {

	control := newControl()
	o.closing = append(o.closing, control)

	c := time.Tick(refresh)

	go func() {
		for {
			select {
			case <-c:
				// log.Println("Check for Trigger")
				if v := o.conn.ReadItem(tag).Value; o.trigger.Fire(v) {
					// log.Println("Triggered.")
					o.Notify()
					o.trigger.Update(v)
				}
			case <-control.close:
				control.done <- true
				return
			}
		}
	}()

	return control
}

//NewObserver returns a new observer
func NewObserver(t Trigger, c Connection) *Observer {
	return &Observer{trigger: t, conn: c, observers: make([]chan bool, 0), closing: make([]*control, 0)}
}
