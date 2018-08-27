package opc

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//OpcMockServerStatic implements an OPC Server that returns the index value plus 1 for each tag.
type OpcMockServerStatic struct {
	Tags []string
}

func (oms *OpcMockServerStatic) Read() dataModel {
	answer := make(dataModel)
	for i, tag := range oms.Tags {
		answer[tag] = float64(i) + 1.0
	}
	return answer
}

//OpcMockServerRandom implements an OPC Server that returns the a random value for each tag.
type OpcMockServerRandom struct {
	Tags []string
}

func (oms *OpcMockServerRandom) Read() dataModel {
	answer := make(dataModel)
	for _, tag := range oms.Tags {
		answer[tag] = rand.Float64()
	}
	return answer
}

func TestOPCDataSync(t *testing.T) {
	odata := NewDataModel()
	running := odata.Sync(&OpcMockServerStatic{[]string{"tag1", "tag2", "tag3"}}, 50*time.Millisecond)
	defer running.Close()

	time.Sleep(200 * time.Millisecond)

	var value float64
	var ok bool

	for i := 0; i < 5; i++ {

		value, ok = odata.Get("tag1")
		if value != 1.0 || !ok {
			t.Fatal("tag1 does not match")
		}

		value, ok = odata.Get("tag2")
		if value != 2.0 || !ok {
			t.Fatal("tag2 does not match")
		}

		value, ok = odata.Get("tag3")
		if value != 3.0 || !ok {
			t.Fatal("tag3 does not match")
		}

		value, ok = odata.Get("tag4")
		if ok {
			t.Fatal("tag4 should not be found")
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func TestOPCDataStop(t *testing.T) {
	odata := NewDataModel()
	running := odata.Sync(&OpcMockServerRandom{[]string{"tag1"}}, 50*time.Millisecond)

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
		running := odata.Sync(&OpcMockServerRandom{[]string{"tag1", "tag2", "tag3"}}, 50*time.Millisecond)
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

/*
func TestOPCDataUpdate(t *testing.T) {
	refreshRate := 50 * time.Millisecond
        odata := NewDataModel()
	running := odata.Sync(&OpcMockServerRandom{[]string{"tag1"}}, refreshRate)

	retries := 10
        sum := 0.0

	for i := -1; i < retries; i++ {
                t := time.Now()
                <-update
                duration := time.Since(t).Seconds()
                if i >= 0 {
                        sum += duration
                }
        }

        log.Println("avg", sum/float64(retries))
        if math.Abs((sum/float64(retries))/1000.0 - 50.0) > 1e-2 {
                t.Fatal("Update interval is too large")
        }

        running.Close()
        ret := <-update
        if ret != struct{}{} {
                t.Fatal("update channel should send nil when channel is closed")
        }

}
*/
