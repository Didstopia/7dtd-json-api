package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Didstopia/7dtd-json-api/server"
)

// API server
type API struct {
	server *server.Server
}

// New returns a pointer to a new API server
func New(server *server.Server) *API {
	api := &API{}
	api.server = server
	return api
}

// Start the API server
func (api *API) Start() error {
	// TODO: Run HTTP server on a separate thread
	// Configure and start the JSON API server
	http.Handle("/", api)
	log.Println("Starting web server on port 8080..")
	if err := http.ListenAndServe(":8080", nil); err != nil { // TODO: Do we need to manually close the server somehow?
		return err
	}
	return nil
}

// Stop the API server
func (api *API) Stop() error {
	return nil
}

func (api *API) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		// json, _ := json.Marshal(api.server.World)
		json, _ := json.MarshalIndent(api.server.World, "", "  ")
		response.WriteHeader(http.StatusOK)
		response.Write(json)
	// case "POST":
	// TODO: Allow sending data to the server as JSON?
	// // Decode the JSON in the body and overwrite 'tom' with it
	// d := json.NewDecoder(request.Body)
	// p := &person{}
	// err := d.Decode(p)
	// if err != nil {
	// 	http.Error(writer, err.Error(), http.StatusInternalServerError)
	// }
	// tom = p
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(response, "Method Not Allowed")
	}
}
