// application implements logic of the game
package application

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Grandeath/Battleship_advanced/application/menu"
	"github.com/Grandeath/Battleship_advanced/connection"
)

// host adress of the api
const (
	host = "https://go-pjatk-server.fly.dev"
)

// StartApp function is starting a game by asking starting condition
// creating base on that starting header. Creates board and client struct
// Pass context to child functions and start RunApp function
func StartApp(ctx context.Context) {
	client := connection.NewClient(host)

	playerIntent := menu.MainMenu(ctx, &client)
	switch playerIntent {
	case menu.ExitTheGame:
		return
	case menu.WaitForChallenge:
		waitForChallenge(ctx, &client)
	case menu.StartGame:
		waitForGameStart(ctx, &client)
	}

	board := NewGuiBoard(true)

	time.Sleep(time.Microsecond * 500)
	desc, err := client.GetDescription(ctx)
	if err != nil {
		log.Println(err)
	}
	board.Description = desc
	time.Sleep(time.Millisecond * 500)
	boardResp, err := client.GetBoard(ctx)
	if err != nil {
		log.Println(err)
	}

	board.CreateBoard(boardResp)
	board.PrintDescription(ctx)

	RunApp(ctx, &board, &client)
}

// RunApp takes context, boardGUI interface and Client interface and run a gane.
// In the for loop function implements running of the programm.
func RunApp(ctx context.Context, board boardGUI, client connection.Client) {
	ctx, cancel := context.WithCancel(ctx)
	// ch receive chosen field from BoardListener
	ch := make(chan string)

	// t blocking channel, which make a loop wait for user chosen field to shot
	t := make(chan struct{})

	// quit channel which close game if value is received form StartBoard goroutine
	quit := make(chan struct{})

	// Start BoardListener in goroutine, which listen user clicks on the board
	go board.BoardListener(ctx, ch, t)

	// Start StartBoard in goroutine to make a board updatable
	go board.StartBoard(ctx, quit)

	board.StartTimer(ctx)

	var fireCoord string

mainloop:
	for {
		// Wait for user turn
		status, err := waitFunction(ctx, client)
		board.SetTurnText("Your turn")
		board.UpdateTImer(status.Timer)
		// Check if sessions was not found. This case happen when user lost becasue of inacrtivity
		if err != nil {
			switch err.(type) {
			case *connection.RequestError:
				if errors.Is(err, &connection.RequestError{StatusCode: 403, Err: "session not found"}) {
					board.LogMessage(err.Error())
					board.LogMessage(status.LastGameStatus)
					break mainloop
				} else {
					board.LogMessage(err.Error())
				}
			default:
				board.LogMessage(err.Error())
			}

		}

		// Check if game anded and log result
		if status.LastGameStatus == "ended" || status.LastGameStatus == "win" || status.LastGameStatus == "lose" || status.GameStatus == "false" {
			board.LogMessage(status.LastGameStatus)
			break mainloop
		}
		// Print enemy shots on user board
		err = board.UpdateYourBoard(status.OppShots)
		if err != nil {
			board.LogMessage(err.Error())
		}
		time.Sleep(time.Millisecond * 500)

		// Wait for channel ch to be ready to pass fire coordinates
		t <- struct{}{}
		select {
		case fireCoord = <-ch:
		case <-quit:
			log.Println("Quit game")
			err = client.DeleteGame(ctx)
			if err != nil {
				log.Println(err)
			}
			break mainloop
		}

		// Shot chosen field to a server and get result of the shot
		resp, err := client.Fire(ctx, fireCoord)
		if err != nil {
			board.LogMessage(err.Error())
		}

		// Print enemy shot to enemy board
		board.FireToBoard(fireCoord, resp)
	}
	cancel()
	time.Sleep(time.Millisecond * 1000)
}

// waitFunction wait for user turn and get status of the game which contain
// if game still going, enemy shots, status of the game for example who won.
func waitFunction(ctx context.Context, client connection.Client) (connection.GameStatus, error) {
	waitingForResponse := true
	gotStatus := connection.GameStatus{}
	for waitingForResponse {
		status, err := client.GetStatus(ctx)
		if err != nil {
			return status, err
		}
		if status.ShouldFire {
			gotStatus = status
			waitingForResponse = false
		}
		if status.LastGameStatus == "ended" || status.LastGameStatus == "win" || status.LastGameStatus == "lose" || status.GameStatus == "false" {
			gotStatus = status
			waitingForResponse = false
		}
		time.Sleep(time.Millisecond * 500)
	}

	return gotStatus, nil
}

func waitForGameStart(ctx context.Context, client connection.Client) {
	err := client.StartGame(ctx)
	if err != nil {
		log.Panic(err)
	}

	waitForGame := true
	for waitForGame {
		status, _ := client.GetStatus(ctx)
		if status.GameStatus == "game_in_progress" {
			waitForGame = false
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func waitForChallenge(ctx context.Context, client connection.Client) {
	err := client.StartGame(ctx)
	if err != nil {
		log.Panic(err)
	}

	var counter int

	waitForGame := true
	for waitForGame {
		status, _ := client.GetStatus(ctx)
		if status.GameStatus == "game_in_progress" {
			waitForGame = false
		}
		time.Sleep(time.Millisecond * 1000)
		if counter == 10 {
			err = client.RefreshWaitingForGame(ctx)
			if err != nil {
				log.Println(err)
			}
			counter = 0
		}
		counter++
	}
}
