package opc

import (
	"reflect"
	"testing"
)

//var client opc.OpcConnection

func TestNewConnectionNoTags(t *testing.T) {
	client := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{},
	)
	client.Close()
}

func TestNewConnectionWithTags(t *testing.T) {
	client := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)
	client.Close()
}

func TestAddTags(t *testing.T) {
	client := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{},
	)
	client.Add("numeric.sin.int64", "numeric.saw.float")
	client.Close()
}

func TestRemoveTags(t *testing.T) {
	client := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)
	client.Remove("numeric.sin.int64")
	client.Remove("numeric.saw.float")
	client.Close()
}

func TestOpcRead(t *testing.T) {
	client := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)

	var item Item

	// should be able to read tag because it has been added
	item = client.ReadItem("numeric.sin.int64")
	if reflect.DeepEqual(item, Item{}) {
		t.Fatal("this test should not have returned an empty item")
	}

	// should be able to read tag because it has been added
	item = client.ReadItem("numeric.fantasy_tag.int64")
	if !reflect.DeepEqual(item, Item{}) {
		t.Fatal("this test should have returned an empty item")
	}

	// read all added tags
	m := client.Read()
	if len(m) != 2 {
		t.Fatal("the map should have only two items")
	}

	client.Close()
}
