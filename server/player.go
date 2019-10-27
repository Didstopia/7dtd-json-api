package server

import (
	"strconv"
	"strings"
)

// Coordinate is a thing
type Coordinate struct {
	X float64 `json:"x,omitempty"`
	Y float64 `json:"y,omitempty"`
	Z float64 `json:"z,omitempty"`
}

func (c *Coordinate) ParseCoordinate(x string, y string, z string) error {
	// Convert each coordinate
	xx, err := strconv.ParseFloat(x, 64)
	yy, err := strconv.ParseFloat(y, 64)
	zz, err := strconv.ParseFloat(z, 64)

	// Return immediately if any errors occurred
	if err != nil {
		return err
	}

	// Store the  coordinates if no errors occurred
	c.X = xx
	c.Y = yy
	c.Z = zz

	// Return nil on success
	return nil
}

// Player is a thing
type Player struct {
	// Index    int         `json:"index"`
	ID       int         `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	Position *Coordinate `json:"position,omitempty"`
	Rotation *Coordinate `json:"rotation,omitempty"`
	// Remote   bool        `json:"remote,omitempty"`
	Health  int    `json:"currentHealth"`
	Deaths  int    `json:"deaths"`
	Zombies int    `json:"zombieKills"`
	Players int    `json:"playerKills"`
	Score   int    `json:"score"`
	Level   int    `json:"level,omitempty"`
	SteamID string `json:"steamId,omitempty"`
	// IP       string      `json:"ip,omitempty"`
	Ping int `json:"ping"`
}

func (p *Player) ParsePlayer(m map[string]string) error {
	// Parse the index
	// index, err := strconv.Atoi(m["index"])
	// if err != nil {
	// 	return err
	// }
	// p.Index = index

	// Parse the ID
	id, err := strconv.Atoi(m["id"])
	if err != nil {
		return err
	}
	p.ID = id

	// Parse the name
	p.Name = m["name"]

	// Parse the position
	pos := strings.Split(m["pos"], ", ")
	p.Position = &Coordinate{}
	if err := p.Position.ParseCoordinate(pos[0], pos[1], pos[2]); err != nil {
		return err
	}

	// Parse the rotation
	rot := strings.Split(m["rot"], ", ")
	p.Rotation = &Coordinate{}
	if err := p.Rotation.ParseCoordinate(rot[0], rot[1], rot[2]); err != nil {
		return err
	}

	// Parse the remote
	// remote, err := strconv.ParseBool(m["remote"])
	// if err != nil {
	// 	return err
	// }
	// p.Remote = remote

	// Parse the health
	health, err := strconv.Atoi(m["health"])
	if err != nil {
		return err
	}
	p.Health = health

	// Parse the deaths
	deaths, err := strconv.Atoi(m["deaths"])
	if err != nil {
		return err
	}
	p.Deaths = deaths

	// Parse the zombies
	zombies, err := strconv.Atoi(m["zombies"])
	if err != nil {
		return err
	}
	p.Zombies = zombies

	// Parse the players
	players, err := strconv.Atoi(m["players"])
	if err != nil {
		return err
	}
	p.Players = players

	// Parse the score
	score, err := strconv.Atoi(m["score"])
	if err != nil {
		return err
	}
	p.Score = score

	// Parse the level
	level, err := strconv.Atoi(m["level"])
	if err != nil {
		return err
	}
	p.Level = level

	// Parse the steamid
	p.SteamID = m["steamid"]

	// Parse the ip
	// p.IP = m["ip"]

	// Parse the ping
	ping, err := strconv.Atoi(m["ping"])
	if err != nil {
		return err
	}
	p.Ping = ping

	return nil
}
