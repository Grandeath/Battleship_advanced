// application implements logic of the game
package application

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Grandeath/Battleship_advanced/application/timer"
	"github.com/Grandeath/Battleship_advanced/connection"
	tl "github.com/grupawp/termloop"
	gui "github.com/grupawp/warships-gui/v2"
)

// GuiBoard contain variables needed to run the board
// and implements function needed to manipulate GUI and make it interactive
// ui - contain user interface
// yourBoard is user board field
// userBoardState - contain state of user's board
// enemyBoard is enemy board field
// enemyBoardState - contain state of enemy's board
// config - contain configuration of the board
// Description contain information about user's and enemy's nick and description
// turnText is field indication whose turn is now
// versusText is field which tell against whome user play
// descText is user description field
// oppDescText is opponent descritpion field
// fireLogText is field which logs your shots
type GuiBoard struct {
	ui                 *gui.GUI
	yourBoard          *gui.Board
	yourBoardState     [10][10]gui.State
	enemyBoard         *gui.Board
	enemyBoardState    [10][10]gui.State
	config             *gui.BoardConfig
	Description        connection.Description
	turnText           *gui.Text
	versusText         *gui.Text
	descText           descriptionField
	oppDescText        descriptionField
	fireLogText        *gui.Text
	legendField        LegendField
	shipLeftCountField ShipLeftCountField
	accuracyField      Accuracy
	timerField         timer.Timer
	enemyCurrentShot   int
}

// NewGuiBoard create GuiBoard
func NewGuiBoard(wantLogs bool) GuiBoard {
	return GuiBoard{ui: gui.NewGUI(wantLogs), config: gui.NewBoardConfig(), enemyCurrentShot: 0}
}

// CreateBoard create board of the game
func (g *GuiBoard) CreateBoard(stateBoard connection.BoardRespons) error {
	g.enemyBoard = gui.NewBoard(50, 6, nil)
	g.ui.Draw(g.enemyBoard)

	states := [10][10]gui.State{}
	for i := range states {
		states[i] = [10]gui.State{}
	}
	g.enemyBoardState = states
	g.enemyBoard.SetStates(g.enemyBoardState)

	for _, ship := range stateBoard.Board {
		column := int(ship[0]) - 65
		row, err := strconv.Atoi(ship[1:])
		if err != nil {
			return err
		}

		states[column][row-1] = gui.Ship
	}
	g.yourBoardState = states
	g.yourBoard = gui.NewBoard(1, 6, nil)
	g.ui.Draw(g.yourBoard)
	g.yourBoard.SetStates(g.yourBoardState)

	return nil
}

// PrintDescription print text fields on the board
func (g *GuiBoard) PrintDescription(ctx context.Context) error {
	g.turnText = gui.NewText(1, 1, "Enemy turn", nil)
	versus := fmt.Sprintf("%s vs %s", g.Description.Nick, g.Description.Opponent)
	g.versusText = gui.NewText(1, 4, versus, nil)
	g.descText = NewDescriptionFieldYour(g.Description.Desc)
	g.oppDescText = NewDescriptionFieldEnemy(g.Description.OppDesc)
	g.fireLogText = gui.NewText(1, 35, "Press on any coordinate to log it.", nil)
	g.accuracyField = NewAccuracyField()
	g.timerField = timer.NewTimer()

	g.legendField = NewLegendField()

	g.shipLeftCountField = NewShipLeftCountField()

	g.ui.Draw(g.turnText)
	g.ui.Draw(g.versusText)
	g.ui.Draw(g.descText.firstLine)
	g.ui.Draw(g.descText.secondLine)
	g.ui.Draw(g.descText.thirdLine)
	g.ui.Draw(g.oppDescText.firstLine)
	g.ui.Draw(g.oppDescText.secondLine)
	g.ui.Draw(g.oppDescText.thirdLine)
	g.ui.Draw(g.fireLogText)
	g.ui.Draw(g.legendField.hitLegend)
	g.ui.Draw(g.legendField.missLegend)
	g.ui.Draw(g.legendField.shipLegend)
	g.ui.Draw(g.shipLeftCountField.FourMastField)
	g.ui.Draw(g.shipLeftCountField.ThreeMastField)
	g.ui.Draw(g.shipLeftCountField.TwoMastField)
	g.ui.Draw(g.shipLeftCountField.OneMastField)
	g.ui.Draw(g.accuracyField.accuracyField)
	g.ui.Draw(g.timerField.ClockField)
	return nil
}

// BoardListener listen to user click
func (g *GuiBoard) BoardListener(ctx context.Context, ch chan<- string, t <-chan struct{}) {
	for {
		<-t
		char := g.enemyBoard.Listen(ctx)
		if len(char) == 0 {
			return
		}
		g.fireLogText.SetText(fmt.Sprintf("Coordinate: %s", char))
		g.ui.Log("Coordinate: %s", char)
		ch <- char
	}
}

// StartBoard starting printing board
func (g *GuiBoard) StartBoard(ctx context.Context, quit chan struct{}) {

	quitKey := tl.Key(tl.KeyCtrlF)

	g.ui.Start(ctx, &quitKey)
	quit <- struct{}{}
}

// Upate enemy board with user shots with its result
func (g *GuiBoard) FireToBoard(coord string, resp connection.FireResponse) error {
	column := int(coord[0]) - 65
	row, err := strconv.Atoi(coord[1:])
	if err != nil {
		return err
	}

	g.accuracyField.ShotNumber++
	switch resp.Result {
	case "hit":
		g.enemyBoardState[column][row-1] = gui.Hit
		g.accuracyField.HitNumber++
	case "miss":
		g.enemyBoardState[column][row-1] = gui.Miss
		g.SetTurnText("Enemy turn")
		g.timerField.TimeStop()
	case "sunk":
		g.ui.Log("Sunk the ship")
		g.enemyBoardState[column][row-1] = gui.Hit
		g.sunkShip(coord)
		g.accuracyField.HitNumber++
	}

	g.ui.Log("Shot at coordinates: %s did %s the target", coord, resp.Result)
	g.accuracyField.updateField()
	g.enemyBoard.SetStates(g.enemyBoardState)
	return nil
}

// LogMessage log text
func (g *GuiBoard) LogMessage(message string) {
	g.ui.Log(message)
}

// UpdateYourBoard update user board with enemy shot with indication if enemy missed or hit a ship
func (g *GuiBoard) UpdateYourBoard(coords []string) error {
	if len(coords) <= g.enemyCurrentShot-1 {
		return nil
	}
	for _, coord := range coords[g.enemyCurrentShot:] {

		g.enemyCurrentShot++
		column := int(coord[0]) - 65
		row, err := strconv.Atoi(coord[1:])
		if err != nil {
			return err
		}

		switch g.yourBoardState[column][row-1] {
		case gui.Ship:
			g.yourBoardState[column][row-1] = gui.Hit
			g.ui.Log("Enemy shot at coordinates: %s did hit target", coord)
		case gui.Hit:
			g.yourBoardState[column][row-1] = gui.Hit
			g.ui.Log("Enemy shot at coordinates: %s did hit target", coord)
		case gui.Empty:
			g.yourBoardState[column][row-1] = gui.Miss
			g.ui.Log("Enemy shot at coordinates: %s did miss target", coord)
		default:
			g.yourBoardState[column][row-1] = gui.Miss
			g.ui.Log("Enemy shot at coordinates: %s did miss target", coord)
		}
		g.yourBoard.SetStates(g.yourBoardState)
	}
	return nil
}

// SetTurnText change text indicating which turn is now
func (g *GuiBoard) SetTurnText(text string) {
	g.turnText.SetText(text)
}

func (g *GuiBoard) sunkShip(coord string) error {
	column := int(coord[0]) - 65
	row, err := strconv.Atoi(coord[1:])
	if err != nil {
		return err
	}
	newNode := NewQuadTree(ShipCoord{row: row - 1, column: column}, g.enemyBoardState, Center)
	g.UpdateShipCountField(len(newNode.GetAllCoords()))
	//check proximity
	g.ui.Log("Destroyed %d mast ship", len(newNode.GetAllCoords()))
	g.HighlighEmptyTiles(row, column)
	return nil
}

func (g *GuiBoard) UpdateShipCountField(shipMastCount int) {
	switch shipMastCount {
	case 4:
		g.shipLeftCountField.FourMastCount = g.shipLeftCountField.FourMastCount - 1
		g.shipLeftCountField.UpdateFourMastCount()
	case 3:
		g.shipLeftCountField.ThreeMastCount = g.shipLeftCountField.ThreeMastCount - 1
		g.shipLeftCountField.UpdateThreeMastCount()
	case 2:
		g.shipLeftCountField.TwoMastCount = g.shipLeftCountField.TwoMastCount - 1
		g.shipLeftCountField.UpdateTwoMastCount()
	case 1:
		g.shipLeftCountField.OneMastCount = g.shipLeftCountField.OneMastCount - 1
		g.shipLeftCountField.UpdateOneMastCount()
	}
}

func (g *GuiBoard) StartTimer(ctx context.Context) {
	g.timerField.StartClock(ctx)
}

func (g *GuiBoard) UpdateTImer(time int) {
	g.timerField.TimeGoes = true
	g.timerField.Time = time
}

func (g *GuiBoard) HighlighEmptyTiles(gotRow int, gotColumn int) {
	stillCircling := true
	heading := North
	column := gotColumn - 1
	row := gotRow - 1

	for {
		if column < 0 {
			break
		}
		if g.enemyBoardState[column][row] == gui.Hit {
			column--
		} else {
			break
		}
	}

	stopColumn := column
	stopRow := row

	for stillCircling {
		if heading == North {
			if row < 10 && row >= 0 {
				if g.enemyBoardState[column+1][row] == gui.Hit {
					if column >= 0 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					if column >= 0 && row-1 >= 0 {
						if g.enemyBoardState[column][row-1] == gui.Hit {
							heading = West
							column--
						} else {
							row--
						}
					} else {
						row--
					}
				} else if g.enemyBoardState[column+1][row] != gui.Hit {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					heading = East
					column++
				}
			} else {
				heading = East
				column++
			}
		}
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}

		if heading == East {
			if column < 10 && column >= 0 {
				if g.enemyBoardState[column][row+1] == gui.Hit {
					if row >= 0 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					if row >= 0 && column+1 < 10 {
						if g.enemyBoardState[column+1][row] == gui.Hit {
							heading = North
							row--
						} else {
							column++
						}
					} else {
						column++
					}
				} else if g.enemyBoardState[column][row+1] != gui.Hit {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					heading = South
					row++
				}
			} else {
				heading = South
				row++
			}
		}
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}
		if heading == South {
			if row < 10 && row >= 0 {
				if g.enemyBoardState[column-1][row] == gui.Hit {
					if column < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					if column < 10 && row+1 < 10 {
						if g.enemyBoardState[column][row+1] == gui.Hit {
							heading = East
							column++
						} else {
							row++
						}
					} else {
						row++
					}
				} else if g.enemyBoardState[column-1][row] != gui.Hit {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					heading = West
					column--
				}
			} else {
				heading = West
				column--
			}
		}
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}
		if heading == West {
			if column < 10 && column >= 0 {
				if g.enemyBoardState[column][row-1] == gui.Hit {
					if row < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					if row < 10 && column-1 >= 0 {
						if g.enemyBoardState[column-1][row] == gui.Hit {
							heading = South
							row++
						} else {
							column--
						}
					} else {
						column--
					}
				} else if g.enemyBoardState[column][row-1] != gui.Hit {
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					heading = North
					row--
				}
			} else {
				heading = North
				row--
			}
		}
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}
	}

}
