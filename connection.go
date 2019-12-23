package opc

import (
	"time"
)

const (
	//OPCDataSource defines constants for Sources when reading data from OPC:
	//Default implementation is OPCCache.
	//From the cache
	OPCCache int32 = 1
	//From the device
	OPCDevice int32 = 2

	//OPCQuality defines the quality of the OPC items:
	//Bad
	OPCQualityBad int16 = 0
	//Good
	OPCQualityGood          int16 = 192
	OPCQualityGoodButForced int16 = 216
	//Maks
	OPCQualityMask int16 = 192
	//Uncertain
	OPCQualityUncertain int16 = 64

	//OPCServerState defines the state of the server:
	//Disconnected
	OPCDisconnected int32 = 6
	//Failed
	OPCFailed int32 = 2
	//Noconfig
	OPCNoconfig int32 = 3
	//Running
	OPCRunning int32 = 1
	//Suspended
	OPCSuspended int32 = 4
	//Test
	OPCTest int32 = 5
)

//Connection represents the interface for the connection to the OPC server.
type Connection interface {
	Add(...string) error
	Remove(string)
	Read() map[string]Item
	ReadItem(string) Item
	Tags() []string
	Write(string, interface{}) error
	Close()
}

//Item stores the result of an OPC item from the OPC server.
type Item struct {
	Value     interface{}
	Quality   int16
	Timestamp time.Time
}

//Good checks the quality of the Item
func (i *Item) Good() bool {
	if i.Quality == OPCQualityGood || i.Quality == OPCQualityGoodButForced {
		return true
	}
	return false
}
