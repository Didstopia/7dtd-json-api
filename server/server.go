package server

import (
	"log"
	"strconv"
	"strings"

	"github.com/Didstopia/7dtd-json-api/server/telnet"
)

const (
	ListingPlayers     telnet.State = "ListingPlayers"
	ListedPlayers      telnet.State = "ListedPlayers"
	GettingTime        telnet.State = "GettingTime"
	GotTime            telnet.State = "GotTime"
	GettingPreferences telnet.State = "GettingPreferences"
	GotPreferences     telnet.State = "GotPreferences"
)

// Server represents the game server
type Server struct {
	World  World `json:"world"`
	telnet *telnet.Telnet
}

// New returns a pointer to a new game server
func New() *Server {
	server := &Server{}
	server.World = World{} // TODO: Update World stuff (players, server info etc.)
	server.World.BloodMoon = &BloodMoon{}
	server.telnet = telnet.New()
	go server.startReceiving()
	return server
}

// Start the game server client
func (server *Server) Start() error {
	return server.telnet.Connect()
}

// Stop the game server client
func (server *Server) Stop() error {
	return server.telnet.Disconnect()
}

// Receive messages from the server
func (server *Server) Receive(message string) {
	log.Println("Server.Receive:", message)

	// Handle getting preferences
	if server.telnet.State == telnet.Idle {
		server.GetPreferences()
	} else if server.telnet.State == GettingPreferences {
		// Parse each preference individually
		if strings.Contains(message, "GamePref.") {
			server.World.ParsePreference(strings.Split(message, "GamePref.")[1])
		}
		// Mark as done from the last known preference
		if strings.Contains(message, "GamePref.ZombiePlayers") {
			server.telnet.SetState(GotPreferences)
		}
	}

	// Handle listing players
	if server.telnet.State == GotPreferences {
		server.ListPlayers()
	} else if server.telnet.State == ListingPlayers {
		// Parse each player individually
		if strings.Contains(message, "steamid") {
			server.World.ParsePlayer(message)
		}
		// Mark as done from the total player listing
		if strings.Contains(message, "Total of ") && strings.Contains(message, " in the game") {
			playerCount, err := strconv.Atoi(strings.Split(strings.Split(message, "Total of ")[1], " in the game")[0])
			if err != nil {
				log.Println("Failed to parse current player count: ", err)
			} else {
				server.World.SetCurrentPlayerCount(playerCount)
			}
			server.telnet.SetState(ListedPlayers)
		}
	}

	// Handle getting time
	if server.telnet.State == ListedPlayers {
		server.GetTime()
	} else if server.telnet.State == GettingTime {
		// Parse the time and mark as done from the known format
		if strings.Contains(message, "Day ") && strings.Contains(message, ", ") && strings.Contains(message, ":") {
			server.World.ParseTime(message)
			server.telnet.SetState(GotTime)
		}
	}

	// Handle moving back to idle
	if server.telnet.State == GotTime {
		server.telnet.SetState(telnet.Idle)
	}
}

// Send commands to the server
func (server *Server) Send(command string) {
	if err := server.telnet.SendCommand(command); err != nil {
		log.Fatal(err)
	}
}

func (server *Server) startReceiving() {
	for {
		// FIXME: Stop receiving if exiting
		server.Receive(<-server.telnet.Messages)
	}
}

// ListPlayers will attempt to send the list players command to the server
func (server *Server) ListPlayers() {
	server.telnet.SetState(ListingPlayers)
	if err := server.telnet.SendCommand("listplayers"); err != nil {
		log.Fatal("Failed to list players: ", err)
	}
}

// GetTime will attempt to send the get time command to the server
func (server *Server) GetTime() {
	server.telnet.SetState(GettingTime)
	if err := server.telnet.SendCommand("gettime"); err != nil {
		log.Fatal("Failed to get time: ", err)
	}
}

// GetPreferences will attempt to send the get game preferences command to the server
func (server *Server) GetPreferences() {
	server.telnet.SetState(GettingPreferences)
	if err := server.telnet.SendCommand("getgamepref"); err != nil {
		log.Fatal("Failed to get time: ", err)
	}
}
