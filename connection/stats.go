package connection

// StatsLeaderboard caontain Stats of 10 best players
type StatsLeaderboard struct {
	GotStats []Stats `json:"stats"`
}

// StatsPlayer contain Stats of the player
type StatsPlayer struct {
	GotStat Stats `json:"stats"`
}

// Stats contain statistick of players
// Games - how many games played
// Nick - nick of player which are those statistics
// Points - points of the player
// Rank - rank of the player
// Wins - how many games player won
type Stats struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}
