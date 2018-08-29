package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/konimarti/opc"
)

// App contains the opc connection and the API routes
type App struct {
	Conn   opc.OpcConnection
	Router *mux.Router
}

// Initialize sets OPC connection and creates routes
func (a *App) Initialize(conn opc.OpcConnection) {
	a.Conn = conn
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/tags", a.getTags).Methods("GET")          // Read
	a.Router.HandleFunc("/tag", a.createTag).Methods("POST")        // Add(...)
	a.Router.HandleFunc("/tag/{id}", a.getTag).Methods("GET")       // ReadItem(id)
	a.Router.HandleFunc("/tag/{id}", a.deleteTag).Methods("DELETE") // Remove(id)
}

// Run starts serving the API
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// getTags returns all tags in the current opc connection, route: /tags
func (a *App) getTags(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, a.Conn.Read())
}

// createTag creates the tags in the opc connection, route: /tag
func (a *App) createTag(w http.ResponseWriter, r *http.Request) {
	var tags []string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tags); err != nil {
		fmt.Println(tags)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := a.Conn.Add(tags...)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Did not add tags")
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{"result": "created"})
}

// getTag returns the opc.Item for the given tag id, route: /tag/{id}
func (a *App) getTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	item := a.Conn.ReadItem(vars["id"])
	empty := opc.Item{}
	if item == empty {
		respondWithError(w, http.StatusNotFound, "tag not found")
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}

// deleteTag removes the tag in the opc connection
func (a *App) deleteTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a.Conn.Remove(vars["id"])
}

// responsWithError is a helper function to return a JSON error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// responsWithJSON is helper function to return the data in JSON encoding
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
