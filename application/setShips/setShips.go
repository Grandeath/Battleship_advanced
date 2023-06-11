// setShips handle a logic to position ships by the player
package setShips

import (
	"context"
	"time"

	"github.com/Grandeath/Battleship_advanced/connection"
)

// StartPositionBoard creates a board and start listening to position ships.
func StartPositionBoard(ctx context.Context, startingHeader *connection.StartingHeader) {
	board := NewSetShipBoard()
	board.CreateBoard()

	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan string)

	// t blocking channel, which make a loop wait for user chosen field to shot.
	t := make(chan struct{})

	// quit channel which close game if value is received form StartBoard goroutine.
	quit := make(chan struct{})

	// Start BoardListener in goroutine, which listen user clicks on the board.
	go board.BoardListener(ctx, ch, t)

	// Start StartBoard in goroutine to make a board updatable.
	go board.StartBoard(ctx, quit)

	var fireCoord string

mainloop:
	for {
		t <- struct{}{}
		select {
		case fireCoord = <-ch:
		case <-quit:
			break mainloop
		}

		board.PlaceShipPart(fireCoord)
	}
	cancel()
	// check if positioned all ships.
	if board.shipToPosition.currentMastCount == 0 {
		startingHeader.Coords = board.ShipPositionSlice
	}

	time.Sleep(time.Millisecond * 200)
}
