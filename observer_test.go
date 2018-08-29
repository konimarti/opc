package opc

import (
	"testing"
	"time"
)

func TestNewObserverOnChange(t *testing.T) {

	refresh := 50 * time.Millisecond
	conn := OpcMockServerRandom{Tags: []string{"tag1", "tag2", "tag3"}}

	// Create Observer
	observer := NewObserver(&OnChange{}, &conn)

	// Start Observer for tag1
	observer.Observe("tag1", refresh)

	// Register with observer
	ch := observer.Channel()

	for i := 0; i < 5; i++ {
		v1 := <-ch
		if v1 != true {
			t.Fatal("this test should work and it should return bool true")
		}
	}
	observer.Close()
}

func TestNewObserverOnValue(t *testing.T) {
	refresh := 50 * time.Millisecond
	conn := OpcMockServerStatic{Tags: []string{"tag1", "tag2", "tag3"}}

	// Create Observer
	observer := NewObserver(&OnValue{1.0}, &conn)
	defer observer.Close()

	// Start Observer for tag1
	observer.Observe("tag1", refresh)

	// Register with observer
	ch := observer.Channel()

	// Start server
	for i := 0; i < 5; i++ {
		v1 := <-ch
		if v1 != true {
			t.Fatal("this test should work and it should return bool true")
		}
	}
}
