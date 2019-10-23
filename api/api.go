package api

import (
	"fmt"
	"log"
	"net/http"
)

// API server
type API struct {
}

// New returns a pointer to a new API server
func New() *API {
	api := &API{}
	return api
}

// Start the API server
func (api *API) Start() error {
	// TODO: Run HTTP server on a separate thread
	// Configure and start the JSON API server
	http.HandleFunc("/", apiHandler)
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

func apiHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		// serverInfo := &ServerInfo{}
		// json, _ := json.Marshal(serverInfo)
		writer.WriteHeader(http.StatusOK)
		// writer.Write(json)
	case "POST":
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
		writer.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(writer, "Method Not Allowed")
	}
}
