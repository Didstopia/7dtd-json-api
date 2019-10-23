package server

// Coordinate is a thing
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// Player is a thing
type Player struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Position Coordinate `json:"position"`
	Rotation Coordinate `json:"rotation"`
	Remote   bool       `json:"remote"`
	Health   int        `json:"health"`
	Deaths   int        `json:"deaths"`
	Zombies  int        `json:"zombies"`
	Players  int        `json:"players"`
	Score    int        `json:"score"`
	Level    int        `json:"level"`
	SteamID  int        `json:"steamid"`
	IP       string     `json:"ip"`
	Ping     int        `json:"ping"`
}
