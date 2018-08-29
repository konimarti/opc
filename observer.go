package opc

import (
	// "log"
	// "errors"
        "io"
	"time"
)

type Trigger interface {
	Fire(interface{}) bool
	Update(interface{})
}

type OnChange struct {
	value interface{}
}

func (t *OnChange) Fire(newValue interface{}) bool { return t.value != newValue }
func (t *OnChange) Update(newValue interface{}) { t.value = newValue }

type OnValue struct {
	value interface{}
}
func (t *OnValue) Fire(newValue interface{}) bool {	return t.value == newValue }
func (t *OnValue) Update(newValue interface{}) {}


type Observer struct {	
	trigger         Trigger
        conn            OpcConnection
        observers       []chan bool
        closing         []*control        
}

func (o *Observer) Close() {
        for _, control := range o.closing {
                control.Close()
        }
}

func (o *Observer) Notify() {
        for _, observer := range o.observers {
                select {
                case <-observer:
                default:
                }
                observer <- true
        }
}

func (o *Observer) Channel() chan bool {
        observer := make(chan bool, 1)
        o.observers = append(o.observers, observer)        
        return observer
}       

func (o *Observer) Observe(tag string, refresh time.Duration) io.Closer {
        
        control := NewControl()
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

func NewObserver(t Trigger, c OpcConnection) *Observer {	
	return &Observer{trigger: t, conn: c, observers: make([]chan bool,0), closing: make([]*control,0)}
}
