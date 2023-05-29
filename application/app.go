package application

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Grandeath/Battleship_advanced/connection"
	tl "github.com/grupawp/termloop"
	gui "github.com/grupawp/warships-gui/v2"
)

type GuiBoard struct {
	ui              *gui.GUI
	yourBoard       *gui.Board
	yourBoardState  [10][10]gui.State
	enemyBoard      *gui.Board
	enemyBoardState [10][10]gui.State
	config          *gui.BoardConfig
	Description     connection.Description
	turnText        *gui.Text
	versusText      *gui.Text
	descText        *gui.Text
	oppDescText     *gui.Text
	fireLogText     *gui.Text
}

func NewGuiBoard(wantLogs bool) GuiBoard {
	return GuiBoard{ui: gui.NewGUI(wantLogs), config: gui.NewBoardConfig()}
}

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

func (g *GuiBoard) PrintDescription(ctx context.Context) error {
	g.turnText = gui.NewText(1, 1, "Your turn", nil)
	versus := fmt.Sprintf("%s vs %s", g.Description.Nick, g.Description.Opponent)
	g.versusText = gui.NewText(1, 4, versus, nil)
	g.descText = gui.NewText(1, 28, g.Description.Desc, nil)
	g.oppDescText = gui.NewText(1, 32, g.Description.OppDesc, nil)
	g.fireLogText = gui.NewText(1, 35, "Press on any coordinate to log it.", nil)

	g.ui.Draw(g.turnText)
	g.ui.Draw(g.versusText)
	g.ui.Draw(g.descText)
	g.ui.Draw(g.oppDescText)
	g.ui.Draw(g.fireLogText)
	return nil
}

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

func (g *GuiBoard) StartBoard() {
	ctx := context.Background()
	quitKey := tl.Key(tl.KeyCtrlF)
	g.ui.Start(ctx, &quitKey)

}

func (g *GuiBoard) FireToBoard(coord string, resp connection.FireResponse) error {
	column := int(coord[0]) - 65
	row, err := strconv.Atoi(coord[1:])
	if err != nil {
		return err
	}

	switch resp.Result {
	case "hit":
		g.enemyBoardState[column][row-1] = gui.Hit
	case "miss":
		g.enemyBoardState[column][row-1] = gui.Miss
	case "Sunk":
		g.enemyBoardState[column][row-1] = gui.Hit
	}

	g.ui.Log("Shot at coordinates: %s did %s target", coord, resp.Result)
	g.enemyBoard.SetStates(g.enemyBoardState)
	return nil
}

func (g *GuiBoard) LogMessage(message string) {
	g.ui.Log(message)
}
