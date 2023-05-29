package application

import (
	"context"

	"github.com/Grandeath/Battleship_advanced/connection"
)

type boardGUI interface {
	CreateBoard(stateBoard connection.BoardRespons) error
	PrintDescription(ctx context.Context) error
	BoardListener(ctx context.Context, ch chan string)
	StartBoard()
}
