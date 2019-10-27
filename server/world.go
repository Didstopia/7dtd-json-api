package server

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type GameTime struct {
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func (time *GameTime) Subtract(t *GameTime) {
	log.Println("Subtract start:", time, t)

	var days uint = 0
	var hours uint = 0
	var minutes uint = 0

	if time.Day > t.Day {
		days = uint(time.Day) - uint(t.Day)
	} else {
		days = uint(t.Day) - uint(time.Day)
	}

	// FIXME: This doesn't take into account 24 hours
	if time.Hour > t.Hour {
		hours = uint(time.Hour) - uint(t.Hour)
	} else {
		hours = uint(t.Hour) - uint(time.Hour)
	}

	// FIXME: This doesn't take into account 60 minutes
	if time.Minute > t.Minute {
		minutes = uint(time.Minute) - uint(t.Minute)
	} else {
		minutes = uint(t.Minute) - uint(time.Minute)
	}

	// TODO: Adjust for negatives values

	log.Println("Subtract result:", days, hours, minutes)

	time.Day = int(days)
	time.Hour = int(hours)
	time.Minute = int(minutes)
}

type BloodMoon struct {
	Frequency int       `json:"frequency"`
	Last      *GameTime `json:"last"`
	Next      *GameTime `json:"next"`
	Countdown *GameTime `json:"countdown"`
}

// World is a thing
type World struct {
	Players            []*Player  `json:"players"`
	CurrentPlayerCount int        `json:"currentPlayerCount"`
	MaxPlayerCount     int        `json:"maxPlayerCount"`
	Time               *GameTime  `json:"time"`
	BloodMoon          *BloodMoon `json:"bloodMoon"`
}

var playerRegex = regexp.MustCompile(`(?P<index>.*?). id=(?P<id>.*?), (?P<name>.*?), pos=\((?P<pos>.*?)\), rot=\((?P<rot>.*?)\), remote=(?P<remote>.*?), health=(?P<health>.*?), deaths=(?P<deaths>.*?), zombies=(?P<zombies>.*?), players=(?P<players>.*?), score=(?P<score>.*?), level=(?P<level>.*?), steamid=(?P<steamid>.*?), ip=(?P<ip>.*?), ping=(?P<ping>.*?)$`)

// ParsePlayers will attempt to parse players from a string
func (world *World) ParsePlayer(player string) bool {
	// Parse values and keys with regex
	values := playerRegex.FindStringSubmatch(player)
	keys := playerRegex.SubexpNames()

	// Create player map
	playerMap := make(map[string]string)
	for i := 1; i < len(keys); i++ {
		playerMap[keys[i]] = values[i]
	}

	// Parse the current player index
	playerIndex, err := strconv.Atoi(playerMap["index"])
	if err != nil {
		log.Println("Failed to parse player: ", err)
		return false
	}

	// Create a player object from the map
	playerObject := &Player{}
	if err := playerObject.ParsePlayer(playerMap); err != nil {
		log.Println("Failed to parse player: ", err)
		return false
	}

	// Initialize empty arrays
	if len(world.Players) == 0 {
		world.Players = make([]*Player, world.MaxPlayerCount)
	}

	// Add or update the player inside the array
	world.Players[playerIndex] = playerObject

	return true
}

func (world *World) ParseTime(time string) bool {
	// Parse the day and hour parts
	dayParts := strings.Split(time, ", ")
	timeParts := strings.Split(dayParts[1], ":")

	// Parse the day
	day, err := strconv.Atoi(strings.Trim(strings.Split(dayParts[0], " ")[1], " "))
	if err != nil {
		log.Println("Failed to parse minutes: ", err)
		return false
	}

	// Parse the hour
	hour, err := strconv.Atoi(strings.Trim(timeParts[0], " "))
	if err != nil {
		log.Println("Failed to parse minutes: ", err)
		return false
	}

	// Parse the minute
	minute, err := strconv.Atoi(strings.Trim(timeParts[1], " "))
	if err != nil {
		log.Println("Failed to parse minutes: ", err)
		return false
	}

	// Store the constructed game time
	world.Time = &GameTime{day, hour, minute}
	// log.Println("Parsed game time:", world.Time)

	// Update the blood moon information
	if world.Time != nil {
		// TODO: Update the last blood moon
		world.BloodMoon.Last = &GameTime{world.Time.Day - world.Time.Day%world.BloodMoon.Frequency, 22, 0}

		// Update the next blood moon
		world.BloodMoon.Next = &GameTime{world.Time.Day + (world.BloodMoon.Frequency - (world.Time.Day % world.BloodMoon.Frequency)), 22, 0}

		// Update the blood moon countdown
		world.BloodMoon.Countdown = &GameTime{world.Time.Day, world.Time.Hour, world.Time.Minute}
		world.BloodMoon.Countdown.Subtract(world.BloodMoon.Next)
	}

	return true
}

func (world *World) ParsePreference(preference string) bool {
	// Parse max player count
	if strings.Contains(preference, "ServerMaxPlayerCount") {
		resultString := strings.Split(preference, "ServerMaxPlayerCount = ")[1]
		resultInt, err := strconv.Atoi(resultString)
		if err != nil {
			log.Println("Failed to parse max. player count: ", err)
			return false
		}
		world.MaxPlayerCount = resultInt
		// log.Println("Parsed max. player count:", world.MaxPlayerCount)
		return true
	}

	// Parse blood moon frequency
	if strings.Contains(preference, "BloodMoonFrequency") {
		resultString := strings.Split(preference, "BloodMoonFrequency = ")[1]
		resultInt, err := strconv.Atoi(resultString)
		if err != nil {
			log.Println("Failed to parse blood moon frequency: ", err)
			return false
		}
		world.BloodMoon.Frequency = resultInt
		// log.Println("Parsed blood moon frequency:", world.BloodMoonFrequency)
		return true
	}

	// Update the next blood moon
	if world.Time != nil {
		// TODO: Also update Last and Countdown
		world.BloodMoon.Next = &GameTime{world.Time.Day + (world.BloodMoon.Frequency - (world.Time.Day % world.BloodMoon.Frequency)), 22, 0}
	}

	return false
}

func (world *World) SetCurrentPlayerCount(playerCount int) {
	// Store the current player count
	world.CurrentPlayerCount = playerCount

	// Resize the player array to match the amount of players online
	if len(world.Players) != world.CurrentPlayerCount {
		newPlayers := make([]*Player, world.CurrentPlayerCount)
		for i := 0; i < len(world.Players); i++ {
			p := world.Players[i]
			if p != nil {
				newPlayers[i] = p
			}
		}
		world.Players = newPlayers
	}

	// Clear out any players past this
	for i := 0; i < len(world.Players); i++ {
		if i > playerCount-1 {
			world.Players[i] = nil
		}
	}

	// Print the players
	log.Println("Listing", world.CurrentPlayerCount, "out of", world.MaxPlayerCount, "players..")
	for index := 0; index < len(world.Players); index++ {
		p := world.Players[index]
		if p != nil {
			log.Println("Player:", world.Players[index])
		}
	}
}
