// connection package handle connection to rest-api server Warships Online API
package connection

import "context"

// Client handle connection to Warships Online API
type Client interface {
	SetStartingHeader(setHeader StartingHeader)
	StartGame(ctx context.Context) error
	GetBoard(ctx context.Context) (BoardRespons, error)
	GetDescription(ctx context.Context) (Description, error)
	Fire(ctx context.Context, coordinates string) (FireResponse, error)
	GetStatus(ctx context.Context) (GameStatus, error)
	GetPlayerList(ctx context.Context) (PlayerList, error)
	DeleteGame(ctx context.Context) error
	GetLeaderBoard(ctx context.Context) (StatsLeaderboard, error)
	GetPlayerScore(ctx context.Context, player string) (StatsPlayer, error)
	RefreshWaitingForGame(ctx context.Context) error
}
