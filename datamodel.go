package opc

import (
	"io"
	"sync"
	"time"
)

//dataModel represents the data structure to hold the OPC data
type dataModel map[string]float64

//data holds the data structure that is refreshed with OPC data.
type data struct {
	tags dataModel
	mu   sync.RWMutex
}

//Get is the thread-safe getter for the tags.
func (d *data) Get(key string) (float64, bool) {
	d.mu.RLock()
	value, ok := d.tags[key]
	d.mu.RUnlock()
	return value, ok
}

//Sync synchronizes the opc server and stores the data into the data model.
func (d *data) Sync(conn OpcConnection, refreshRate time.Duration) io.Closer {

	close := make(chan bool)
	done := make(chan bool)

	ticker := time.NewTicker(refreshRate)

	go func() {
		for {
			select {
			case <-ticker.C:
				update := conn.Read()
				d.mu.Lock()
				for key, value := range update {
					d.tags[key] = value
				}
				d.mu.Unlock()
			case <-close:
				ticker.Stop()
				done <- true
				return
			}
		}
	}()

	return &control{close, done}
}

//NewDataModel returns an OPC Data struct.
func NewDataModel() data {
	return data{tags: make(dataModel)}
}

type control struct {
	close chan bool
	done  chan bool
}

func (c *control) Close() error {
	if c.close != nil && c.done != nil {
		c.close <- true
		<-c.done
	}
	return nil
}

func newControl() *control {
	return nil
}
