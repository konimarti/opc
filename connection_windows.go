// +build windows

package opc

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	OleInit()
}

//OleInit initializes OLE.
func OleInit() {
	ole.CoInitializeEx(0, 0)
	// log.Println("OLE initialized.")
}

//OleRealease realeses OLE resources in opcAutomation.
func OleRelease() {
	ole.CoUninitialize()
	// log.Println("OLE released.")
}

//AutomationObject loads the OPC Automation Wrapper and handles to connection to the OPC Server.
type AutomationObject struct {
	unknown *ole.IUnknown
	object  *ole.IDispatch
}

//Connect establishes a connection to the OPC Server on node.
//It returns a reference to AutomationItems and error message.
func (ao *AutomationObject) Connect(server string, node string) (*AutomationItems, error) {

	// check if server is running, if yes then disconnect
	if ao.IsConnected() {
		_, err := oleutil.CallMethod(ao.object, "Disconnect")
		if err != nil {
			log.Println("Failed to disconnect. Trying to connect anyway..")
		}
	}

	// try to connect to opc server and check for error
	log.Printf("Connecting to %s on node %s\n", server, node)
	_, err := oleutil.CallMethod(ao.object, "Connect", server, node)
	if err != nil {
		log.Println("Connection failed.")
		return nil, errors.New("Connection failed")
	}

	// set up opc groups and items
	opcGroups, err := oleutil.GetProperty(ao.object, "OPCGroups")
	if err != nil {
		//log.Println(err)
		return nil, errors.New("Cannot get OPCGroups property.")
	}
	opcGrp, err := oleutil.CallMethod(opcGroups.ToIDispatch(), "Add")
	if err != nil {
		// log.Println(err)
		return nil, errors.New("Cannot add new OPC Group.")
	}
	addItemObject, err := oleutil.GetProperty(opcGrp.ToIDispatch(), "OPCItems")
	if err != nil {
		// log.Println(err)
		return nil, errors.New("Cannot get OPC Items.")
	}

	opcGroups.ToIDispatch().Release()
	opcGrp.ToIDispatch().Release()

	log.Println("Connected.")

	return NewAutomationItems(addItemObject.ToIDispatch()), nil
}

//TryConnect loops over the nodes array and tries to connect to any of the servers.
func (ao *AutomationObject) TryConnect(server string, nodes []string) (*AutomationItems, error) {
	var errResult string
	for _, node := range nodes {
		items, err := ao.Connect(server, node)
		if err == nil {
			return items, err
		}
		errResult = errResult + err.Error() + "\n"
	}
	return nil, errors.New("TryConnect was not successfull: " + errResult)
}

//IsConnected check if the server is properly connected and up and running.
func (ao *AutomationObject) IsConnected() bool {
	if ao.object == nil {
		return false
	}
	state_vt, err := oleutil.GetProperty(ao.object, "ServerState")
	if err != nil {
		log.Println("GetProperty call for ServerState failed", err)
		return false
	}
	if state_vt.Value().(int32) != OPCRunning {
		return false
	}
	return true
}

//Close releases the OLE objects in the AutomationObject.
func (ao *AutomationObject) Close() {
	ao.object.Release()
	ao.unknown.Release()
}

func NewAutomationObject() *AutomationObject {
	unknown, err := oleutil.CreateObject("OPC.Automation.1")
	if err != nil {
		log.Println("Could not load OPC Automation object")
		return &AutomationObject{}
	}
	opc, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Println("Could not QueryInterface")
		return &AutomationObject{}
	}
	object := AutomationObject{
		unknown: unknown,
		object:  opc,
	}
	return &object
}

//AutomationItems store the OPCItems from OPCGroup and does the bookkeeping
//for the individual OPC items. Tags can added, removed, and read.
type AutomationItems struct {
	addItemObject *ole.IDispatch
	items         map[string]*ole.IDispatch
}

//addSingle adds the tag and returns an error. Client handles are not implemented yet.
func (ai *AutomationItems) addSingle(tag string) error {
	clientHandle := int32(1)
	item, err := oleutil.CallMethod(ai.addItemObject, "AddItem", tag, clientHandle)
	if err != nil {
		return errors.New(tag + ":" + err.Error())
	}
	ai.items[tag] = item.ToIDispatch()
	return nil
}

//Add accepts a variadic parameters of tags.
func (ai *AutomationItems) Add(tags ...string) error {
	var errResult string
	for _, tag := range tags {
		err := ai.addSingle(tag)
		if err != nil {
			errResult = err.Error() + errResult
		}
	}
	if errResult == "" {
		return nil
	}
	return errors.New(errResult)
}

//Remove removes the tag.
func (ai *AutomationItems) Remove(tag string) {
	item, ok := ai.items[tag]
	if ok {
		item.Release()
	}
	delete(ai.items, tag)
}

//readFromOPC reads from the server and returns an Item and error.
func (ai *AutomationItems) readFromOpc(opcitem *ole.IDispatch) (Item, error) {
	v := ole.NewVariant(ole.VT_R4, 0)
	q := ole.NewVariant(ole.VT_INT, 0)
	ts := ole.NewVariant(ole.VT_DATE, 0)

	//read tag from opc server and monitor duration in seconds
	t := time.Now()
	_, err := oleutil.CallMethod(opcitem, "Read", OPCCache, &v, &q, &ts)
	opcReadsDuration.Observe(time.Since(t).Seconds())

	if err != nil {
		opcReadsCounter.WithLabelValues("failed").Inc()
		return Item{}, err
	}
	opcReadsCounter.WithLabelValues("success").Inc()

	return Item{
		Value:     v.Value(),
		Quality:   q.Value().(int16),
		Timestamp: ts.Value().(time.Time),
	}, nil
}

//writeToOPC writes value to opc tag and return an error
func (ai *AutomationItems) writeToOpc(opcitem *ole.IDispatch, value interface{}) error {
	_, err := oleutil.CallMethod(opcitem, "Write", value)
	if err != nil {
		// TODO: Prometheus Monitoring
		//opcWritesCounter.WithLabelValues("failed").Inc()
		return err
	}
	//opcWritesCounter.WithLabelValues("failed").Inc()
	return nil
}

//Close closes the OLE objects in AutomationItems.
func (ai *AutomationItems) Close() {
	for key, opcitem := range ai.items {
		opcitem.Release()
		delete(ai.items, key)
	}
	ai.addItemObject.Release()
}

//NewAutomationItems returns a new AutomationItems instance.
func NewAutomationItems(opcitems *ole.IDispatch) *AutomationItems {
	ai := AutomationItems{addItemObject: opcitems, items: make(map[string]*ole.IDispatch)}
	return &ai
}

//opcRealServer implements the OpcConnection interface.
//It has the AutomationObject embedded for connecting to the server
//and an AutomationItems to facilitate the OPC items bookkeeping.
type opcConnectionImpl struct {
	*AutomationObject
	*AutomationItems
	Server string
	Nodes  []string
	Tags   []string
	mu     sync.Mutex
}

//ReadItem returns an Item for a specific tag.
func (conn *opcConnectionImpl) ReadItem(tag string) Item {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	opcitem, ok := conn.AutomationItems.items[tag]
	if ok {
		item, err := conn.AutomationItems.readFromOpc(opcitem)
		if err == nil {
			return item
		}
		log.Printf("Cannot read %s: %s. Trying to fix.", tag, err)
		conn.fix()
	} else {
		log.Printf("Tag %s not found. Add it first before reading it.", tag)
	}
	return Item{}
}

//Write writes a value to the OPC Server.
func (conn *opcConnectionImpl) Write(tag string, value interface{}) error {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	opcitem, ok := conn.AutomationItems.items[tag]
	if ok {
		return conn.AutomationItems.writeToOpc(opcitem, value)
	} else {
		log.Printf("Tag %s not found. Add it first before writing to it.", tag)
	}
	return errors.New("No Write performed")
}

//Read returns a map of the values of all added tags.
func (conn *opcConnectionImpl) Read() map[string]interface{} {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	allValues := map[string]interface{}{}
	for tag, opcitem := range conn.AutomationItems.items {
		item, err := conn.AutomationItems.readFromOpc(opcitem)
		if err != nil {
			log.Printf("Cannot read %s: %s. Trying to fix.", tag, err)
			conn.fix()
			break
		}
		allValues[tag] = item.Value
	}
	return allValues
}

//fix tries to reconnect if connection is lost by creating a new connection
//with AutomationObject and creating a new AutomationItems instance.
func (conn *opcConnectionImpl) fix() {
	var err error
	if !conn.IsConnected() {
		for {
			conn.AutomationItems.Close()
			conn.AutomationItems, err = conn.TryConnect(conn.Server, conn.Nodes)
			if err != nil {
				log.Println(err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if conn.Add(conn.Tags...) == nil {
				log.Printf("Added %d tags", len(conn.Tags))
			}
			break
		}
	}
}

//Close closes the embedded types.
func (conn *opcConnectionImpl) Close() {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	if conn.AutomationObject != nil {
		conn.AutomationObject.Close()
	}
	if conn.AutomationItems != nil {
		conn.AutomationObject.Close()
	}
}

//NewConnection establishes a connection to the OpcServer object.
func NewConnection(server string, nodes []string, tags []string) OpcConnection {
	object := NewAutomationObject()
	items, err := object.TryConnect(server, nodes)
	if err != nil {
		panic(err)
	}
	err = items.Add(tags...)
	if err != nil {
		panic(err)
	}
	conn := opcConnectionImpl{
		AutomationObject: object,
		AutomationItems:  items,
		Server:           server,
		Nodes:            nodes,
		Tags:             tags,
	}

	return &conn
}
