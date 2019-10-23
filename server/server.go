package server

import "github.com/Didstopia/7dtd-json-api/server/telnet"

// Server represents the game server
type Server struct {
	World  World `json:"world"`
	telnet *telnet.Telnet
}

// New returns a pointer to a new game server
func New() *Server {
	server := &Server{}
	server.World = World{} // TODO: Update World stuff (players, server info etc.)
	server.telnet = telnet.New()
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
