package setships

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Grandeath/Battleship_advanced/application/ship"
	tl "github.com/grupawp/termloop"
	gui "github.com/grupawp/warships-gui/v2"
)

type SetShipBoard struct {
	ui                *gui.GUI
	setBoard          *gui.Board
	setBoardState     [10][10]gui.State
	config            *gui.BoardConfig
	shipToPosition    PositionField
	placementCount    int
	ShipPositionSlice []string
	placementBucket   []ship.ShipCoord
}

func NewSetShipBoard() SetShipBoard {
	return SetShipBoard{ui: gui.NewGUI(true), config: gui.NewBoardConfig()}
}

func (s *SetShipBoard) CreateBoard() {
	s.setBoard = gui.NewBoard(1, 1, nil)

	states := [10][10]gui.State{}
	for i := range states {
		states[i] = [10]gui.State{}
	}

	s.setBoardState = states
	s.setBoard.SetStates(s.setBoardState)

	s.shipToPosition = NewPositionField()

	s.ui.Draw(s.setBoard)
	s.ui.Draw(s.shipToPosition.currentShipFIeld)
}

func (s *SetShipBoard) BoardListener(ctx context.Context, ch chan<- string, t <-chan struct{}) {
	for {
		<-t
		char := s.setBoard.Listen(ctx)
		if len(char) == 0 {
			return
		}
		s.ui.Log("Coordinate: %s", char)
		ch <- char
	}
}

func (s *SetShipBoard) StartBoard(ctx context.Context, quit chan struct{}) {
	quitKey := tl.Key(tl.KeyCtrlF)
	s.ui.Start(ctx, &quitKey)
	quit <- struct{}{}
}

func (s *SetShipBoard) PlaceShipPart(coord string) error {
	column := int(coord[0]) - 65
	row, err := strconv.Atoi(coord[1:])

	if err != nil {
		return err
	}

	if s.setBoardState[column][row-1] != gui.Miss && s.setBoardState[column][row-1] != gui.Ship && s.shipToPosition.currentMastCount != 0 {
		s.setBoardState[column][row-1] = gui.Ship
		s.ui.Log("Placed ship at %s", coord)
		s.placementCount++
		s.placementBucket = append(s.placementBucket, ship.ShipCoord{Row: row, Column: column})
	}

	if s.placementCount == s.shipToPosition.currentMastCount {
		s.CheckIfCorrectPosition(row, column)
		s.placementCount = 0
	}
	s.setBoard.SetStates(s.setBoardState)
	return nil
}

func (s *SetShipBoard) CheckIfCorrectPosition(row int, column int) {
	newNode := ship.NewQuadTree(ship.ShipCoord{Row: row - 1, Column: column}, s.setBoardState, ship.Center, gui.Ship)

	currentShipCoords := newNode.GetAllCoords()

	if s.shipToPosition.currentMastCount == len(currentShipCoords) {
		for _, coord := range currentShipCoords {
			s.ShipPositionSlice = append(s.ShipPositionSlice, transfromCoord(coord.Row, coord.Column))
		}
		s.HighlighEmptyTiles(row, column)
		s.shipToPosition.NextShip()
		s.placementBucket = s.placementBucket[:0]
	} else {
		for _, coord := range s.placementBucket {
			s.setBoardState[coord.Column][coord.Row-1] = gui.Empty
		}
		s.placementBucket = s.placementBucket[:0]
	}

}

func transfromCoord(row int, coord int) string {
	columnLetter := rune(coord + 65)
	return fmt.Sprintf("%c%d", columnLetter, row+1)
}

func (s *SetShipBoard) HighlighEmptyTiles(gotRow int, gotColumn int) {
	heading := ship.North
	column := gotColumn - 1
	row := gotRow - 1

	for {
		if column < 0 {
			break
		}
		if s.setBoardState[column][row] == gui.Ship {
			column--
		} else {
			break
		}
	}

	stopColumn := column
	stopRow := row

mainloop:
	for {
		if heading == ship.North {
			if row < 10 && row >= 0 {
				if s.setBoardState[column+1][row] == gui.Ship {
					if column >= 0 {
						s.setBoardState[column][row] = gui.Miss
					}
					if column >= 0 && row-1 >= 0 {
						if s.setBoardState[column][row-1] == gui.Ship {
							heading = ship.West
							column--
						} else {
							row--
						}
					} else {
						row--
					}
				} else if s.setBoardState[column+1][row] != gui.Ship {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						s.setBoardState[column][row] = gui.Miss
					}
					heading = ship.East
					column++
				}
			} else {
				heading = ship.East
				column++
			}
		}
		if row == stopRow && column == stopColumn {
			break mainloop
		}

		if heading == ship.East {
			if column < 10 && column >= 0 {
				if s.setBoardState[column][row+1] == gui.Ship {
					if row >= 0 {
						s.setBoardState[column][row] = gui.Miss
					}
					if row >= 0 && column+1 < 10 {
						if s.setBoardState[column+1][row] == gui.Ship {
							heading = ship.North
							row--
						} else {
							column++
						}
					} else {
						column++
					}
				} else if s.setBoardState[column][row+1] != gui.Ship {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						s.setBoardState[column][row] = gui.Miss
					}
					heading = ship.South
					row++
				}
			} else {
				heading = ship.South
				row++
			}
		}
		if row == stopRow && column == stopColumn {
			break mainloop
		}
		if heading == ship.South {
			if row < 10 && row >= 0 {
				if s.setBoardState[column-1][row] == gui.Ship {
					if column < 10 {
						s.setBoardState[column][row] = gui.Miss
					}
					if column < 10 && row+1 < 10 {
						if s.setBoardState[column][row+1] == gui.Ship {
							heading = ship.East
							column++
						} else {
							row++
						}
					} else {
						row++
					}
				} else if s.setBoardState[column-1][row] != gui.Ship {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						s.setBoardState[column][row] = gui.Miss
					}
					heading = ship.West
					column--
				}
			} else {
				heading = ship.West
				column--
			}
		}
		if row == stopRow && column == stopColumn {
			break mainloop
		}
		if heading == ship.West {
			if column < 10 && column >= 0 {
				if s.setBoardState[column][row-1] == gui.Ship {
					if row < 10 {
						s.setBoardState[column][row] = gui.Miss
					}
					if row < 10 && column-1 >= 0 {
						if s.setBoardState[column-1][row] == gui.Ship {
							heading = ship.South
							row++
						} else {
							column--
						}
					} else {
						column--
					}
				} else if s.setBoardState[column][row-1] != gui.Ship {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						s.setBoardState[column][row] = gui.Miss
					}
					heading = ship.North
					row--
				}
			} else {
				heading = ship.North
				row--
			}
		}
		if row == stopRow && column == stopColumn {
			break mainloop
		}
	}

}
