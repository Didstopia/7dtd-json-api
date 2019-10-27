package main

import (
	"log"

	"github.com/Didstopia/7dtd-json-api/api"
	"github.com/Didstopia/7dtd-json-api/server"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create a new server client and connect to the server
	server := server.New()
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start the server client: ", err)
	}

	// Create a new API client and start it
	api := api.New(server)
	if err := api.Start(); err != nil {
		log.Fatal("Failed to start API server: ", err)
	}
}
