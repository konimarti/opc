package opc

import (
	"reflect"
	"testing"
)

func TestOPCBrowser(t *testing.T) {
	browser, err := CreateBrowser(
		"Graybox.Simulator",
		[]string{"localhost"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if browser.Name != "root" {
		t.Fatal("structure of browser tree is compromised: root")
	}
	if browser.Branches[0].Name != "options" {
		t.Fatal("structure of browser tree is compromised: options")
	}
	if len(browser.Branches[0].Leaves) != 4 {
		t.Fatal("structure of browser tree is compromised: number of leaves for options")
	}
}

func TestNewConnectionNoTags(t *testing.T) {
	client, _ := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{},
	)
	client.Close()
}

func TestNewConnectionWithTags(t *testing.T) {
	client,_ := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)
	client.Close()
}

func TestNewConnectionWrongServer(t *testing.T) {
	client, err := NewConnection(
		"Graybox.Simulator.NOTREAL",
		[]string{"localhost"},
		[]string{},
	)
	client.Close()

	if err == nil {
		t.Fatal("this test should return an error because server does not exist")
	}

}

func TestNewConnectionWrongNode(t *testing.T) {
	client, err := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost.NOTREAL"},
		[]string{},
	)
	client.Close()

	if err == nil {
		t.Fatal("this test should return an error because node does not exist")
	}

}


func TestAddTags(t *testing.T) {
	client, _ := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{},
	)
	defer client.Close()
	client.Add("numeric.sin.int64", "numeric.saw.float")
}

func TestRemoveTags(t *testing.T) {
	client, _ := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)
	defer client.Close()
	client.Remove("numeric.sin.int64")
	client.Remove("numeric.saw.float")
}

func TestOpcRead(t *testing.T) {
	client, _ := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)
	defer client.Close()

	var item Item

	// should be able to read tag because it has been added
	item = client.ReadItem("numeric.sin.int64")
	if reflect.DeepEqual(item, Item{}) {
		t.Fatal("this test should not have returned an empty item")
	}

	// should not be able to read tag because it does not exist
	item = client.ReadItem("numeric.fantasy_tag.int64")
	if !reflect.DeepEqual(item, Item{}) {
		t.Fatal("this test should have returned an empty item")
	}

	// read all added tags (items)
	m := client.Read()
	if len(m) != 2 {
		t.Fatal("the map should have only two items")
	}

	// check Good() of Item
	if item.Quality == OPCQualityGood {
		if item.Good() != true {
			t.Fatal("failed to check quality of item")
		}
	} else {
		if item.Good() == true {
			t.Fatal("failed to check quality of item")
		}
	}
}

func TestOpcWrite(t *testing.T) {
	client, _ := NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)
	defer client.Close()

	var config = []struct {
		Tag     string
		Payload interface{}
		Want    interface{}
	}{
		{
			Tag:     "storage.numeric.reg01",
			Payload: 0.12,
			Want:    0.12,
		},
		{
			Tag:     "storage.numeric.reg02",
			Payload: 2,
			Want:    2.0,
		},
		{
			Tag:     "storage.string.reg01",
			Payload: "Hello",
			Want:    "Hello",
		},
		{
			Tag:     "storage.bool.reg01",
			Payload: true,
			Want:    true,
		},
	}

	for _, cfg := range config {

		// write new frequency to non-existing tag which should fail
		err := client.Write(cfg.Tag, cfg.Payload)
		if err == nil {
			t.Fatal("this test should fail because tag has not been added yet and cannot be written to")
		}

		// add tag
		client.Add(cfg.Tag)

		// write new frequency to existing tag which should succeed
		err = client.Write(cfg.Tag, cfg.Payload)
		if err != nil {
			t.Fatal("this test should not fail because new value should be written to tag")
		}

		// read tag and check if value has been changed
		item := client.ReadItem(cfg.Tag)
		if item.Value != cfg.Want {
			t.Fatalf("tag has not been set to value. Got %v but expected %v", item.Value, cfg.Want)
		}
	}
}
