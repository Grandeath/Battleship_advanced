package connection

type StatsLeaderboard struct {
	GotStats []Stats `json:"stats"`
}

type StatsPlayer struct {
	GotStat Stats `json:"stats"`
}

type Stats struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}
