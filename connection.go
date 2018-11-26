package opc

import (
	"time"
)

const (
	//OPCDataSource defines constants for Sources when reading data from OPC.
	//Default implementation is OPCCache.
	oPCCache  int32 = 1
	oPCDevice int32 = 2

	//OPCQuality defines constants for OPCItems that were read.
	oPCQualityBad       int16 = 0
	oPCQualityGood      int16 = 192
	oPCQualityMask      int16 = 192
	oPCQualityUncertain int16 = 64

	//OPCServerState defines the state of the server.
	oPCDisconnected int32 = 6
	oPCFailed       int32 = 2
	oPCNoconfig     int32 = 3
	oPCRunning      int32 = 1
	oPCSuspended    int32 = 4
	oPCTest         int32 = 5

	//OPCErrors defines errors when reading OPCItems.
	oPCBadRights = -1073479674
	oPCBadType   = -1073479676
	//..
)

//Connection represents the interface for the connection to the OPC server.
type Connection interface {
	Add(...string) error
	Remove(string)
	Read() map[string]interface{}
	ReadItem(string) Item
	Write(string, interface{}) error
	Close()
}

//Item stores the result of an OPC item from the OPC server.
type Item struct {
	Value     interface{}
	Quality   int16
	Timestamp time.Time
}
