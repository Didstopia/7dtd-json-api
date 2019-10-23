package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

// ServerInfo is the top level object,
// containing all server information
type ServerInfo struct {
	Players int `json:"players"`
}

// Server interface
type Server struct {
	hostname    string
	port        string
	password    string
	connection  *telnet.Conn
	lastMessage string
}

func (server *Server) Write(b []byte) (n int, err error) {
	// Parse the incoming data as as string (technically it's always a single byte/character)
	dataString := string(b)

	// fmt.Printf("RECEIVED: %#v\n", dataString)

	// Bail out if the string is empty
	if len(dataString) <= 0 {
		return len(b), errors.New("Invalid or zero length data")
	}

	// Detect the end of line (null terminator, new line or return)
	endOfLine := dataString == "\x00" || dataString == "\n" || dataString == "\r"

	// TODO: Figure out a way to more easily parse incoming data etc.

	// TODO: Execute "getgamepref" and parse "GamePref.BloodMoonFrequency = 7",
	//       then combine that with "gettime" to easily calculate the horde day/ETA to horde

	// TODO: Figure out how to exit after calling exit (or after the telnet connection ends?)

	// Handle message building
	if endOfLine {
		// Guard against empty strings
		if len(server.lastMessage) > 0 {
			fmt.Printf("Message received: %#v\n", server.lastMessage)
			server.Authenticate()
			server.ListPlayers()
			server.GetTime()
			server.Exit()
			server.lastMessage = ""
		}
	} else {
		server.lastMessage += dataString
	}

	// Return the original length and nil error on success
	return len(b), nil
}

// Authenticate checks for a password request and sends the password
func (server *Server) Authenticate() {
	if err := server.sendCommandIfContains(server.password, "Please enter password:"); err != nil {
		log.Fatal("Failed to authenticate: ", err)
	}
}

// ListPlayers will attempt to send the list players command to the server
func (server *Server) ListPlayers() {
	if err := server.sendCommandIfContains("listplayers", "Press 'exit' to end session."); err != nil {
		log.Fatal("Failed to list players: ", err)
	}
}

// GetTime will attempt to send the get time command to the server
func (server *Server) GetTime() {
	if err := server.sendCommandIfContains("gettime", "Press 'exit' to end session."); err != nil {
		log.Fatal("Failed to get time: ", err)
	}
}

// Exit will send the exit command to the server
func (server *Server) Exit() {
	if err := server.sendCommandIfContains("exit", "Executing command 'gettime'"); err != nil {
		log.Fatal("Failed to exit: ", err)
	}
}

func (server *Server) sendCommandIfContains(command string, contains string) error {
	if strings.Contains(server.lastMessage, contains) {
		return server.sendCommand(command)
	}
	return nil
}

func (server *Server) sendCommand(command string) error {
	log.Println("Sending command", "'"+command+"'"+" to the server..")
	p := []byte(command + "\r")
	if _, err := oi.LongWrite(server.connection, p); err != nil {
		return err
	}
	return nil
}

func main() {
	// Create a new server object
	server := &Server{}

	// Verify that we have a valid server
	server.hostname = os.Getenv("SERVER")
	if len(server.hostname) <= 0 {
		log.Fatal("Missing required environment variable: SERVER")
	}

	// Verify that we have a valid port
	server.port = os.Getenv("PORT")
	if len(server.port) <= 0 {
		log.Fatal("Missing required environment variable: port")
	}

	// Load the password but don't verify it, as it's not required
	server.password = os.Getenv("PASSWORD")

	// Create a new telnet client
	// telnetClient := telnet.StandardCaller

	// reader, _ := Server.(Reader)
	// writer, _ := Server.(Writer)

	// go telnetClient.CallTELNET(telnet.NewContext(), serverObject, serverObject)

	// TODO: Run telnet client on a separate thread
	// Attempt to connect to the telnet server
	connectionString := server.hostname + ":" + server.port
	log.Println("Connecting to " + connectionString + "..")
	connection, err := telnet.DialTo(connectionString)
	// err := telnet.DialToAndCall(connectionString, telnetClient)
	if err != nil {
		log.Fatal("Failed to connect to telnet server: ", err)
	}
	server.connection = connection
	defer connection.Close()
	log.Println("Connection open!")

	// TODO: Read data from the telnet server
	go func(writer io.Writer, reader io.Reader) {
		var buffer [1]byte
		p := buffer[:]
		for {
			n, err := reader.Read(p)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}
			oi.LongWrite(writer, p)
		}
	}(server, connection)

	// TODO: Run HTTP server on a separate thread
	// Configure and start the JSON API server
	http.HandleFunc("/", apiHandler)
	log.Println("Starting web server on port 8080..")
	if err := http.ListenAndServe(":8080", nil); err != nil { // TODO: Do we need to manually close the server somehow?
		log.Fatal("Failed to start web server: ", err)
	}
}

func apiHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		serverInfo := &ServerInfo{}
		json, _ := json.Marshal(serverInfo)
		writer.WriteHeader(http.StatusOK)
		writer.Write(json)
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
