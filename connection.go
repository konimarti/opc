package opc

import (
	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	//OPCDataSource defines constants for Sources when reading data from OPC.
	OPCCache  int32 = 1
	OPCDevice int32 = 2

	//OPCQuality defines constants for OPCItems that were read.
	OPCQualityBad       int32 = 0
	OPCQualityGood      int32 = 192
	OPCQualityMask      int32 = 192
	OPCQualityUncertain int32 = 64

	//OPCServerState defines the state of the server.
	OPCDisconnected int32 = 6
	OPCFailed       int32 = 2
	OPCNoconfig     int32 = 3
	OPCRunning      int32 = 1
	OPCSuspended    int32 = 4
	OPCTest         int32 = 5

	//OPCErrors defines errors when reading OPCItems
	OPCBadRights = -1073479674
	OPCBadType   = -1073479676
	//..
)

var (
	//opcAutomation is a non-exported, global variable containing the OPC Automation Wrapper.
	//It needs to be initialized with opc.Ole_init() and realsed with opc.Ole_release().
	//If not initialized, it is nil.
	opcAutomation *ole.IDispatch
	once          sync.Once
)

//GetOpcAutomation implements the singleton pattern and returns a pointer to the automation object.
func GetOpcAutomation() *ole.IDispatch {
	once.Do(func() {
		//CoInitializeEx is necessary to get concurrent COM model.
		ole.CoInitializeEx(0, 0)
		unknown, err := oleutil.CreateObject("OPC.Automation.1")
		if err != nil {
			log.Panicln("Could not load OPC Automation object")
		}
		opc, err := unknown.QueryInterface(ole.IID_IDispatch)
		if err != nil {
			log.Panicln("Could not QueryInterface")
		}
		log.Println("OLE initialized.")

		opcAutomation = opc
	})
	return opcAutomation
}

//OleInit initializes OLE resources in opcAutomation.
func OleInit() {
	//Nothing to be done. Using lazy initialization of OLE.
}

//OleRealease realeses OLE resources in opcAutomation.
func OleRelease() {
	if opcAutomation != nil {
		opc_groups := oleutil.MustGetProperty(opcAutomation, "OPCGroups").ToIDispatch()
		oleutil.MustCallMethod(opc_groups, "RemoveAll")
		oleutil.MustCallMethod(opcAutomation, "Disconnect")
		opcAutomation.Release()
	}
	ole.CoUninitialize()
	log.Println("OLE released.")
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

//OpcConnection represents the interface for OPC Servers to be used with opc.Data.
type OpcConnection interface {
	Read() dataModel
}

//opcRealServer implements an actual OLE object over DCOM and is not exported.
type opcConnectionImpl struct {
	Tags          []string
	Group         *ole.IDispatch
	Items         []*ole.IDispatch
	ClientHandles []int32
	ServerHandles []int32
	logger        *log.Logger
}

//Read reads all tags from OPC Server with a OpcGroup.SyncRead call which is blocking.
func (ors *opcConnectionImpl) Read() dataModel {

	answer := make(dataModel)

	v := ole.NewVariant(ole.VT_R4, 0)
	q := ole.NewVariant(ole.VT_INT, 0)
	ts := ole.NewVariant(ole.VT_DATE, 0)

	for i, item := range ors.Items {
                //read tag from opc server and monitor duration in seconds
		t := time.Now()
		_, err := oleutil.CallMethod(item, "Read", OPCCache, &v, &q, &ts)
		opcReadsDuration.Observe(time.Since(t).Seconds())

		if err != nil {
			ors.logger.Printf("Could not read tag (%s). Try next.\n", ors.Tags[i])
			opcReadsCounter.WithLabelValues("failed").Inc()
			continue
		} else {
			opcReadsCounter.WithLabelValues("success").Inc()
		}

		// ors.logger.Printf("Tag (%s) Value(%v) Quality(%v) Timestamp(%v)\n", ors.Tags[i], v.Value(), q.Value(), ts.Value())

		var value float64
		switch v.Value().(type) {
		case int32:
			value = float64(v.Value().(int32))
		case int64:
			value = float64(v.Value().(int64))
		case float32:
			value = float64(v.Value().(float32))
		case float64:
			value = float64(v.Value().(float64))
		default:
			ors.logger.Println("Could not type assert return value for tag (%s). Try next.", ors.Tags[i])
			continue
		}

		answer[ors.Tags[i]] = value
	}

	return answer
}

//NewOpcConnection establishes a connection to the OpcServer object.
func NewOpcConnection(server string, node string, tags []string, output io.Writer) OpcConnection {

	logger := log.New(output, "OPCConn ", log.LstdFlags)

	state := oleutil.MustGetProperty(GetOpcAutomation(), "ServerState").Value().(int32)
	logger.Println("Server State: ", state)

	if state != OPCRunning {
		logger.Printf("Connecting with server (%s) on node (%s)\n", server, node)
		oleutil.MustCallMethod(GetOpcAutomation(), "Connect", server, node).ToIDispatch()
	}

	state = oleutil.MustGetProperty(GetOpcAutomation(), "ServerState").Value().(int32)
	logger.Println("Server State: ", state)

	if state != OPCRunning {
		logger.Println("OPC Server is not running Returning nil. Code: ", state)
		return nil
	}
	logger.Println("Server is up.")

	opc_groups := oleutil.MustGetProperty(GetOpcAutomation(), "OPCGroups").ToIDispatch()
	opc_group := oleutil.MustCallMethod(opc_groups, "Add").ToIDispatch()
	opc_items := oleutil.MustGetProperty(opc_group, "OPCItems").ToIDispatch()

	clienthandles := make([]int32, len(tags))
	for i, _ := range tags {
		clienthandles[i] = int32(i + 1)
	}

	items := make([]*ole.IDispatch, len(tags))
	serverhandles := make([]int32, len(tags))
	for i, tag := range tags {
		items[i] = oleutil.MustCallMethod(opc_items, "AddItem", tag, clienthandles[i]).ToIDispatch()
		serverhandles[i] = oleutil.MustGetProperty(items[i], "ServerHandle").Value().(int32)
		logger.Printf("Tag (%s) created with ClientHandle (%d) ServerHandle (%d).\n", tag, clienthandles[i], serverhandles[i])
	}

	return &opcConnectionImpl{Tags: tags, Group: opc_group, Items: items, ClientHandles: clienthandles, ServerHandles: serverhandles, logger: logger}
}
