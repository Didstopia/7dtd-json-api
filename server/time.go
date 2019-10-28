package server

// GameTime represents the in-game time
type GameTime struct {
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}
