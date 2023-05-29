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
	ui          *gui.GUI
	yourBoard   *gui.Board
	enemyBoard  *gui.Board
	config      *gui.BoardConfig
	Description connection.Description
	turnText    *gui.Text
	versusText  *gui.Text
	descText    *gui.Text
	oppDescText *gui.Text
	fireLogText *gui.Text
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

	g.enemyBoard.SetStates(states)

	for _, ship := range stateBoard.Board {
		column := int(ship[0]) - 65
		row, err := strconv.Atoi(ship[1:])
		if err != nil {
			return err
		}

		states[column][row-1] = gui.Ship
	}

	g.yourBoard = gui.NewBoard(1, 6, nil)
	g.ui.Draw(g.yourBoard)
	g.yourBoard.SetStates(states)

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

func (g *GuiBoard) BoardListener(ctx context.Context, ch chan string) {
	for {
		char := g.enemyBoard.Listen(ctx)
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
