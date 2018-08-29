package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/konimarti/opc"
	"github.com/konimarti/opc/api"
)

var a api.App

func TestMain(m *testing.M) {
	a = api.App{}

	client := opc.NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{},
	)
	defer client.Close()

	a = api.App{}
	a.Initialize(client)

	// run "main"
	code := m.Run()

	os.Exit(code)
}

// test empty tags, route: /tags
func TestEmptyTags(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tags", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "{}" {
		t.Errorf("Expected an empty map. Got %s", body)
	}
}

// test return when tag is non-existent, route: /tag/{id}
func TestNonExistingProduct(t *testing.T) {
	req, _ := http.NewRequest("GET", "/tag/non.existing.tag", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "tag not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'tag not found'. Got '%s'", m["error"])
	}
}

// test create tags, route: /tag
func TestCreateTag(t *testing.T) {
	payload := []byte(`["numeric.sin.float","numeric.sin.int32","numeric.saw.float"]`)

	req, _ := http.NewRequest("POST", "/tag", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["result"] != "created" {
		t.Errorf("Expected result to be 'created'. Got '%v'", m["result"])
	}
}

// test read an item, route: /tag/{id}
func TestGetTag(t *testing.T) {
	// add tag "numeric.sin.float"
	req, _ := http.NewRequest("POST", "/tag", bytes.NewBuffer([]byte(`["numeric.sin.float"]`)))
	executeRequest(req)

	req, _ = http.NewRequest("GET", "/tag/numeric.sin.float", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var item opc.Item
	json.Unmarshal(response.Body.Bytes(), &item)
	empty := opc.Item{}
	if item == empty {
		t.Errorf("Expected item with value, quality, and timestamp. Got %v", item)
	}

	if item.Quality == 0 {
		t.Errorf("Quality code not ok. Got %v", item.Quality)
	}
}

// helper functions
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response %d. Got %d", expected, actual)
	}
}
