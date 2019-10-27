package telnet

import (
	"errors"
	"io"
	"log"
	"os"
	"strings"

	"github.com/reiver/go-oi"
	tel "github.com/reiver/go-telnet"
)

type State string

const (
	Connecting     State = "Connecting"
	Connected      State = "Connected"
	Authenticating State = "Authenticating"
	Authenticated  State = "Authenticated"
	Idle           State = "Idle"
	Exiting        State = "Exiting"
)

// Telnet client specific to the game server
type Telnet struct {
	hostname    string
	port        string
	password    string
	connection  *tel.Conn
	lastMessage string

	Messages chan string
	State    State
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

	// Create a channel for the incoming messages
	telnet.Messages = make(chan string)

	return telnet
}

// Connect to the telnet server
func (telnet *Telnet) Connect() error {
	telnet.SetState(Connecting)

	// TODO: Run telnet client on a separate thread
	// Attempt to connect to the telnet server
	connectionString := telnet.hostname + ":" + telnet.port
	connection, err := tel.DialTo(connectionString)
	if err != nil {
		log.Fatal("Failed to connect to telnet server: ", err)
	}
	telnet.connection = connection
	telnet.SetState(Connected)

	// Read data from the telnet server
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
	// Ignore any data i exiting
	if telnet.State == Exiting {
		return
	}

	// Parse the incoming data as as string (technically it's always a single byte/character)
	dataString := string(b)

	// fmt.Printf("RECEIVED: %#v\n", dataString)

	// Bail out if the string is empty
	if len(dataString) <= 0 {
		return len(b), errors.New("Invalid or zero length data")
	}

	// Detect the end of line (null terminator, new line or return)
	endOfLine := dataString == "\x00" || dataString == "\n" || dataString == "\r"

	// Handle message building
	if endOfLine {
		// Guard against empty strings
		if len(telnet.lastMessage) > 0 {
			// fmt.Printf("Message received: %#v\n", telnet.lastMessage)

			// Handle connecting
			if telnet.State == Connecting {
				if strings.Contains(telnet.lastMessage, "Connected with 7DTD server.") {
					telnet.SetState(Connected)
				}
			}

			// Handle authentication
			if telnet.State == Connected {
				if strings.Contains(telnet.lastMessage, "Please enter password:") {
					telnet.Authenticate()
				}
			} else if telnet.State == Authenticating {
				if strings.Contains(telnet.lastMessage, "Logon successful.") {
					telnet.SetState(Authenticated)
				}
			} else if telnet.State == Authenticated {
				telnet.SetState(Idle)
			}

			// TODO: Execute "getgamepref" and parse "GamePref.BloodMoonFrequency = 7",
			//       then combine that with "gettime" to easily calculate the horde day/ETA to horde

			// telnet.Exit()
			telnet.Messages <- telnet.lastMessage
			telnet.lastMessage = ""
		}
	} else {
		telnet.lastMessage += dataString
	}

	// Return the original length and nil error on success
	return len(b), nil
}

func (telnet *Telnet) SetState(state State) {
	log.Println("Telnet state changing from", telnet.State, "to", state)
	telnet.State = state
}

// Authenticate checks for a password request and sends the password
func (telnet *Telnet) Authenticate() {
	telnet.SetState(Authenticating)
	if err := telnet.SendCommand(telnet.password); err != nil {
		log.Fatal("Failed to authenticate: ", err)
	}
}

// Exit will send the exit command to the server
func (telnet *Telnet) Exit() {
	telnet.SetState(Exiting)
	if err := telnet.SendCommand("exit"); err != nil {
		log.Fatal("Failed to exit: ", err)
	}

	// TODO: Obviously don't do this here
	os.Exit(0)
}

// SendCommand sends a command to the server
func (telnet *Telnet) SendCommand(command string) error {
	// log.Println("Sending command", "'"+command+"'"+" to the server..")
	p := []byte(command + "\r")
	if _, err := oi.LongWrite(telnet.connection, p); err != nil {
		return err
	}
	return nil
}
