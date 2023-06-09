package connection

// Description store information from GetDescription function about your and enemy login and description
type Description struct {
	Desc     string `json:"desc"`
	Nick     string `json:"nick"`
	OppDesc  string `json:"opp_desc"`
	Opponent string `json:"opponent"`
}

// GameStatus store information get from GetStatus function about
// GameStatus - status of the game
// ShoulFire - boolean if this is your turn to fire
// OppShots - oponnent shots coordinates during their turn
// LastGameStatus - information about outcome of last game
// Timer - information about your how long turn will last
type GameStatus struct {
	GameStatus     string   `json:"game_status"`
	ShouldFire     bool     `json:"should_fire"`
	OppShots       []string `json:"opp_shots"`
	LastGameStatus string   `json:"last_game_status"`
	Timer          int      `json:"timer"`
}
