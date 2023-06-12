package connection

// StartingHeader header for StartGame function which store information about game you want to begin
// Desc - if specified send your description to a server if not server will assign it automatically for you
// Nick - if specified send your Nick to a server if not server will assign it automatically for you
// TargetNick - if specified attack chosen player
// Wpbot - if true you chose to fighr WPbot
// Coords - contain positions of the ships position by the player
type StartingHeader struct {
	Desc       string   `json:"desc"`
	Nick       string   `json:"nick"`
	TargetNick string   `json:"target_nick"`
	Wpbot      bool     `json:"wpBot"`
	Coords     []string `json:"coords"`
}

// PlayerList decode information got from GetPlayerList
// GameStatus - tell the status of the game of the player
// Nick - tell a nick of the waiting player
type PlayerList []struct {
	GameStatus string `json:"game_status"`
	Nick       string `json:"nick"`
}

// BoardRespons decode information get from GetBoard function about how the board of the game looks
type BoardRespons struct {
	Board []string
}
