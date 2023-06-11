// application implements logic of the game
package application

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Grandeath/Battleship_advanced/application/ship"
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
	quitField          *gui.Text
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
	g.descText = NewDescriptionFieldYours(g.Description.Desc)
	g.oppDescText = NewDescriptionFieldEnemy(g.Description.OppDesc)
	g.fireLogText = gui.NewText(1, 35, "Press on any coordinate to log it.", nil)
	g.accuracyField = NewAccuracyField()
	g.timerField = timer.NewTimer()
	g.quitField = gui.NewText(30, 1, "To quit the game press ctrl+f", nil)

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
	g.ui.Draw(g.quitField)
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

// sunkShip when shot sunk a ship this function taking coordinates of last shot and update ShipCountField
// and highlight empty tiles with HighlighEmptyTiles function around ship
func (g *GuiBoard) sunkShip(coord string) error {
	// convert string coordinates to number of row and column
	column := int(coord[0]) - 65
	row, err := strconv.Atoi(coord[1:])
	if err != nil {
		return err
	}

	// create ship tree which store all coordinates of this ship
	newNode := ship.NewQuadTree(ship.ShipCoord{Row: row - 1, Column: column}, g.enemyBoardState, ship.Center, gui.Hit)

	// update ship count with lenghth of the ship
	g.UpdateShipCountField(len(newNode.GetAllCoords()))
	//check proximity
	g.ui.Log("Destroyed %d mast ship", len(newNode.GetAllCoords()))

	// sate tiles around the sunken ship to miss
	g.HighlighEmptyTiles(row, column)
	return nil
}

// UpdateShipCountField update shipLeftCountField on the mast count number of sunken ship
func (g *GuiBoard) UpdateShipCountField(shipMastCount int) {
	switch shipMastCount {
	case 4:
		g.shipLeftCountField.FourMastCount = g.shipLeftCountField.FourMastCount - 1
		g.shipLeftCountField.UpdateFourMastField()
	case 3:
		g.shipLeftCountField.ThreeMastCount = g.shipLeftCountField.ThreeMastCount - 1
		g.shipLeftCountField.UpdateThreeMastField()
	case 2:
		g.shipLeftCountField.TwoMastCount = g.shipLeftCountField.TwoMastCount - 1
		g.shipLeftCountField.UpdateTwoMastField()
	case 1:
		g.shipLeftCountField.OneMastCount = g.shipLeftCountField.OneMastCount - 1
		g.shipLeftCountField.UpdateOneMastField()
	}
}

// StartTimer starting timerField which listen how much time left player have to decide on shot
func (g *GuiBoard) StartTimer(ctx context.Context) {
	g.timerField.StartClock(ctx)
}

// UpdateTimer take time and set counter to this number
func (g *GuiBoard) UpdateTImer(time int) {
	g.timerField.TimeGoes = true
	g.timerField.Time = time
}

// HighlighEmptyTiles algorithm which set fields around the ship to miss and takes row and column coordinates of last shot.
// Script move clocwise.
func (g *GuiBoard) HighlighEmptyTiles(gotRow int, gotColumn int) {
	// if still circling around the ship
	stillCircling := true

	// moving direction of the scrip. In this case algorithm start at heading north
	heading := ship.North
	column := gotColumn - 1
	row := gotRow - 1

	// moving starting column at first empty left tile to the last shot
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

	// Set stopping coordinates
	stopColumn := column
	stopRow := row

	// Starting algoritm
	for stillCircling {
		// check if moving direction is to North
		if heading == ship.North {
			// check if row is in the range of enemyBoardState slice. If not change a heading to East
			// when movinb North You will always change a row number so if row is out of bounds that mean script need to change direction
			// I decided that move is clocwise
			if row < 10 && row >= 0 {
				// Check if on the right site is ship
				if g.enemyBoardState[column+1][row] == gui.Hit {
					// check if column is out of bound and set the field to miss if script is on the board
					if column >= 0 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					// check if when moving script forward is still on board
					if column >= 0 && row-1 >= 0 {
						// check if ahead is ship. If is change heading to West and rise a column. If not just head north
						if g.enemyBoardState[column][row-1] == gui.Hit {
							heading = ship.West
							column--
						} else {
							row--
						}
					} else {
						row--
					}
					// If on the right side there is no ship
				} else if g.enemyBoardState[column+1][row] != gui.Hit {
					// set current field to miss if it is not out of bounds and change moving direction to East. This part make a corners to miss and turn to right.
					if column >= 0 && column < 10 && row >= 0 && row < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					heading = ship.East
					column++
				}
			} else {
				heading = ship.East
				column++
			}
		}
		// check if script made full circle
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}
		// rest of the script is as North but with appropriate corrections
		if heading == ship.East {
			if column < 10 && column >= 0 {
				if g.enemyBoardState[column][row+1] == gui.Hit {
					if row >= 0 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					if row >= 0 && column+1 < 10 {
						if g.enemyBoardState[column+1][row] == gui.Hit {
							heading = ship.North
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
					heading = ship.South
					row++
				}
			} else {
				heading = ship.South
				row++
			}
		}
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}
		if heading == ship.South {
			if row < 10 && row >= 0 {
				if g.enemyBoardState[column-1][row] == gui.Hit {
					if column < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					if column < 10 && row+1 < 10 {
						if g.enemyBoardState[column][row+1] == gui.Hit {
							heading = ship.East
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
					heading = ship.West
					column--
				}
			} else {
				heading = ship.West
				column--
			}
		}
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}
		if heading == ship.West {
			if column < 10 && column >= 0 {
				if g.enemyBoardState[column][row-1] == gui.Hit {
					if row < 10 {
						g.enemyBoardState[column][row] = gui.Miss
					}
					if row < 10 && column-1 >= 0 {
						if g.enemyBoardState[column-1][row] == gui.Hit {
							heading = ship.South
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
					heading = ship.North
					row--
				}
			} else {
				heading = ship.North
				row--
			}
		}
		if row == stopRow && column == stopColumn {
			stillCircling = false
		}
	}

}
