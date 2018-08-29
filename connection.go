package opc

import (
	"time"
)

const (
	//OPCDataSource defines constants for Sources when reading data from OPC.
	OPCCache  int32 = 1
	OPCDevice int32 = 2

	//OPCQuality defines constants for OPCItems that were read.
	OPCQualityBad       int16 = 0
	OPCQualityGood      int16 = 192
	OPCQualityMask      int16 = 192
	OPCQualityUncertain int16 = 64

	//OPCServerState defines the state of the server.
	OPCDisconnected int32 = 6
	OPCFailed       int32 = 2
	OPCNoconfig     int32 = 3
	OPCRunning      int32 = 1
	OPCSuspended    int32 = 4
	OPCTest         int32 = 5

	//OPCErrors defines errors when reading OPCItems.
	OPCBadRights = -1073479674
	OPCBadType   = -1073479676
	//..
)

//OpcConnection represents the interface for the connection to the OPC server.
type OpcConnection interface {
	Add(...string) error
	Remove(string)
	Read() map[string]interface{}
	ReadItem(string) Item
	Close()
}

//Item stores the result of an OPC item from the OPC server.
type Item struct {
	Value     interface{}
	Quality   int16
	Timestamp time.Time
}
