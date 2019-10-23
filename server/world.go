package server

// World is a thing
type World struct {
	Time    Time
	Players []Player `json:"players"`
}

// ParsePlayers will attempt to parse players from a string
func (world *World) ParsePlayers(players string) error {
	// TODO: Parse players from the string: https://scene-si.org/2017/09/02/parsing-strings-with-go/

	// TODO: Each player is sent as its own message, as seen below
	// TODO: Player list is done when the "Total of X in the game" message is sent
	// 0. id=229, Rumilus, pos=(-1394.0, 59.1, 611.0), rot=(-26.7, 181.4, 0.0), remote=True, health=150, deaths=7, zombies=1414, players=0, score=1351, level=135, steamid=76561198008931473, ip=84.249.70.57, ping=18
	// 1. id=236, Tepaya, pos=(-1397.0, 59.1, 611.5), rot=(-19.7, 11674.7, 0.0), remote=True, health=170, deaths=5, zombies=272, players=0, score=222, level=111, steamid=76561198079774759, ip=84.249.68.209, ping=17
	// Total of 2 in the game

	return nil
}
