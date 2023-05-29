package application

import (
	"context"

	"github.com/Grandeath/Battleship_advanced/connection"
)

type boardGUI interface {
	CreateBoard(stateBoard connection.BoardRespons) error
	PrintDescription(ctx context.Context) error
	BoardListener(ctx context.Context, ch chan<- string, t <-chan struct{})
	StartBoard()
	FireToBoard(coord string, resp connection.FireResponse) error
	LogMessage(message string)
}
