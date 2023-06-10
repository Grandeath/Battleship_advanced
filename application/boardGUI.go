// application implements logic of the game
package application

import (
	"context"

	"github.com/Grandeath/Battleship_advanced/connection"
)

// boardGUI implements GUI of the game
type boardGUI interface {
	CreateBoard(stateBoard connection.BoardRespons) error
	PrintDescription(ctx context.Context) error
	BoardListener(ctx context.Context, ch chan<- string, t <-chan struct{})
	StartBoard(ctx context.Context, quit chan struct{})
	FireToBoard(coord string, resp connection.FireResponse) error
	LogMessage(message string)
	UpdateYourBoard(coords []string) error
	SetTurnText(text string)
	StartTimer(ctx context.Context)
	UpdateTImer(time int)
}
