package opc

import (
	// "log"
	// "fmt"

	"testing"
	"time"
)

func TestOPCDataSync(t *testing.T) {
	odata := NewDataModel()
	running := odata.Sync(&OpcMockServerStatic{Tags: []string{"tag1", "tag2", "tag3"}}, 50*time.Millisecond)
	defer running.Close()

	time.Sleep(200 * time.Millisecond)

	var value interface{}
	var ok bool

	for i := 0; i < 5; i++ {

		value, ok = odata.Get("tag1")
		if value.(float64) != 1.0 || !ok {
			t.Fatal("tag1 does not match")
		}

		value, ok = odata.Get("tag2")
		if value.(float64) != 2.0 || !ok {
			t.Fatal("tag2 does not match")
		}

		value, ok = odata.Get("tag3")
		if value.(float64) != 3.0 || !ok {
			t.Fatal("tag3 does not match")
		}

		_, ok = odata.Get("tag4")
		if ok {
			t.Fatal("tag4 should not be found")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TestOPCDataStop(t *testing.T) {
	odata := NewDataModel()
	running := odata.Sync(&OpcMockServerRandom{Tags: []string{"tag1"}}, 50*time.Millisecond)

	time.Sleep(200 * time.Millisecond)

	value1, _ := odata.Get("tag1")

	time.Sleep(200 * time.Millisecond)

	value2, _ := odata.Get("tag1")

	if value1 == value2 {
		t.Fatal("values should be different")
	}

	running.Close()

	time.Sleep(200 * time.Millisecond)

	value1, _ = odata.Get("tag1")

	time.Sleep(200 * time.Millisecond)

	value2, _ = odata.Get("tag1")

	if value1 != value2 {
		t.Fatal("values should be the same if no more updating")
	}
}

func TestOPCDataStopAfterRecord(t *testing.T) {

	c := make(chan bool)

	go func() {
		odata := NewDataModel()
		running := odata.Sync(&OpcMockServerRandom{Tags: []string{"tag1", "tag2", "tag3"}}, 50*time.Millisecond)
		running.Close()
		c <- true
	}()

	for {
		select {
		case <-c:
			return
		case <-time.After(2 * time.Second):
			t.Fatal("time out while closing after syncronizing")
		}
	}
}
