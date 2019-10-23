package telnet

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/reiver/go-oi"
	tel "github.com/reiver/go-telnet"
)

// Telnet client specific to the game server
type Telnet struct {
	hostname    string
	port        string
	password    string
	connection  *tel.Conn
	lastMessage string
}

// New returns a pointer to a telnet client
func New() *Telnet {
	telnet := &Telnet{}

	// Verify that we have a valid server
	telnet.hostname = os.Getenv("SERVER")
	if len(telnet.hostname) <= 0 {
		log.Fatal("Missing required environment variable: SERVER")
	}

	// Verify that we have a valid port
	telnet.port = os.Getenv("PORT")
	if len(telnet.port) <= 0 {
		log.Fatal("Missing required environment variable: port")
	}

	// Load the password but don't verify it, as it's not required
	telnet.password = os.Getenv("PASSWORD")

	return telnet
}

// Connect to the telnet server
func (telnet *Telnet) Connect() error {
	// TODO: Run telnet client on a separate thread
	// Attempt to connect to the telnet server
	connectionString := telnet.hostname + ":" + telnet.port
	log.Println("Connecting to " + connectionString + "..")
	connection, err := tel.DialTo(connectionString)
	// err := telnet.DialToAndCall(connectionString, telnetClient)
	if err != nil {
		log.Fatal("Failed to connect to telnet server: ", err)
	}
	telnet.connection = connection
	// defer connection.Close()
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
	}(telnet, connection)

	return nil
}

// Disconnect from the telnet server
func (telnet *Telnet) Disconnect() error {
	return telnet.connection.Close()
}

func (telnet *Telnet) Write(b []byte) (n int, err error) {
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
		if len(telnet.lastMessage) > 0 {
			fmt.Printf("Message received: %#v\n", telnet.lastMessage)
			telnet.Authenticate()
			telnet.ListPlayers()
			telnet.GetTime()
			// telnet.Exit()
			telnet.lastMessage = ""
		}
	} else {
		telnet.lastMessage += dataString
	}

	// Return the original length and nil error on success
	return len(b), nil
}

// Authenticate checks for a password request and sends the password
func (telnet *Telnet) Authenticate() {
	if err := telnet.sendCommandIfContains(telnet.password, "Please enter password:"); err != nil {
		log.Fatal("Failed to authenticate: ", err)
	}
}

// ListPlayers will attempt to send the list players command to the server
func (telnet *Telnet) ListPlayers() {
	if err := telnet.sendCommandIfContains("listplayers", "Press 'exit' to end session."); err != nil {
		log.Fatal("Failed to list players: ", err)
	}
}

// GetTime will attempt to send the get time command to the server
func (telnet *Telnet) GetTime() {
	if err := telnet.sendCommandIfContains("gettime", "Press 'exit' to end session."); err != nil {
		log.Fatal("Failed to get time: ", err)
	}
}

// Exit will send the exit command to the server
func (telnet *Telnet) Exit() {
	if err := telnet.sendCommandIfContains("exit", "Executing command 'gettime'"); err != nil {
		log.Fatal("Failed to exit: ", err)
	}
}

func (telnet *Telnet) sendCommandIfContains(command string, contains string) error {
	if strings.Contains(telnet.lastMessage, contains) {
		return telnet.sendCommand(command)
	}
	return nil
}

func (telnet *Telnet) sendCommand(command string) error {
	log.Println("Sending command", "'"+command+"'"+" to the server..")
	p := []byte(command + "\r")
	if _, err := oi.LongWrite(telnet.connection, p); err != nil {
		return err
	}
	return nil
}
