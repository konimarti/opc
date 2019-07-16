package opc

import (
	// "log"
	// "fmt"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//Embedded type for mock servers
type emptyServer struct{}

func (es *emptyServer) Add(...string) error             { return nil }
func (es *emptyServer) Remove(string)                   {}
func (es *emptyServer) Write(string, interface{}) error { return nil }
func (es *emptyServer) Close()                          {}
func (es *emptyServer) Tags() []string                  { return []string{} }

//OpcMockServerStatic implements an OPC Server that returns the index value plus 1 for each tag.
type OpcMockServerStatic struct {
	*emptyServer
	TagList []string
}

func (oms *OpcMockServerStatic) ReadItem(tag string) Item {
	items := oms.Read()
	return items[tag]
}

func (oms *OpcMockServerStatic) Read() map[string]Item {
	answer := make(map[string]Item)
	for i, tag := range oms.TagList {
		answer[tag] = Item{float64(i) + 1.0, OPCQualityGood, time.Now()}
	}
	return answer
}

//OpcMockServerRandom implements an OPC Server that returns a random value for each tag.
type OpcMockServerRandom struct {
	*emptyServer
	TagList []string
}

func (oms *OpcMockServerRandom) ReadItem(tag string) Item {
	items := oms.Read()
	return items[tag]
}

func (oms *OpcMockServerRandom) Read() map[string]Item {
	answer := make(map[string]Item)
	for _, tag := range oms.TagList {
		answer[tag] = Item{rand.Float64(), OPCQualityGood, time.Now()}
	}
	return answer
}

//OpcMockServerWakeUp implements an OPC Server that returns 1.0 for a certain duration then a random value for each tag.
type OpcMockServerWakeUp struct {
	*emptyServer
	TagList []string
	AtSleep bool
	mu      sync.Mutex
}

func (oms *OpcMockServerWakeUp) WakeUpAfter(sleep time.Duration) {
	oms.AtSleep = true
	go func() {
		time.Sleep(sleep)
		oms.mu.Lock()
		oms.AtSleep = false
		oms.mu.Unlock()
	}()
}

func (oms *OpcMockServerWakeUp) ReadItem(tag string) Item {
	items := oms.Read()
	return items[tag]
}

func (oms *OpcMockServerWakeUp) Read() map[string]Item {
	answer := make(map[string]Item)

	oms.mu.Lock()
	defer oms.mu.Unlock()

	for _, tag := range oms.TagList {
		if oms.AtSleep {
			answer[tag] = Item{1.0, OPCQualityGood, time.Now()}
		} else {
			answer[tag] = Item{rand.Float64(), OPCQualityGood, time.Now()}
		}
	}

	return answer
}

//FallAsleep Server, sets to 2.0 after time period (opposite of WakeUp server)
type OpcMockServerFallAsleep struct {
	*emptyServer
	TagList []string
	AtSleep bool
	mu      sync.Mutex
}

func (oms *OpcMockServerFallAsleep) FallAsleepAfter(sleep time.Duration) {
	oms.AtSleep = false
	go func() {
		time.Sleep(sleep)
		oms.mu.Lock()
		oms.AtSleep = true
		oms.mu.Unlock()
	}()
}

func (oms *OpcMockServerFallAsleep) ReadItem(tag string) Item {
	items := oms.Read()
	return items[tag]
}

func (oms *OpcMockServerFallAsleep) Read() map[string]Item {
	answer := make(map[string]Item)

	oms.mu.Lock()
	defer oms.mu.Unlock()

	for _, tag := range oms.TagList {
		if oms.AtSleep {
			answer[tag] = Item{2.0, OPCQualityGood, time.Now()}
		} else {
			answer[tag] = Item{rand.Float64(), OPCQualityGood, time.Now()}
		}
	}

	return answer
}
