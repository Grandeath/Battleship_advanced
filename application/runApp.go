package application

import (
	"context"

	"log"
	"time"

	"github.com/Grandeath/Battleship_advanced/connection"
)

const (
	host = "https://go-pjatk-server.fly.dev"
)

func StartApp() {
	client := connection.NewClient(host)

	board := NewGuiBoard(true)
	client.StartingHeader = connection.StartingHeader{Wpbot: true}

	ctx := context.Background()
	err := client.StartGame(ctx)
	if err != nil {
		log.Println(err)
	}

	waitForGame := true
	for waitForGame {
		status, _ := client.GetStatus(ctx)
		if status.GameStatus == "game_in_progress" {
			waitForGame = false
		}
		time.Sleep(time.Millisecond * 1000)
	}

	desc, err := client.GetDescription(ctx)
	if err != nil {
		log.Println(err)
	}
	board.Description = desc

	boardResp, err := client.GetBoard(ctx)
	if err != nil {
		log.Println(err)
	}

	board.CreateBoard(boardResp)
	board.PrintDescription(ctx)

	RunApp(ctx, &board, &client)
}

func RunApp(ctx context.Context, board boardGUI, client connection.Client) {
	// gameStillGoing := true
	ch := make(chan string)

	go board.BoardListener(ctx, ch)
	board.StartBoard()
	// fmt.Println(<-ch)

	// for gameStillGoing {
	// 	status, err := waitFunction(ctx, client)
	// 	if errors.Is(err, fmt.Errorf("unexpected status code: 403")) {
	// 		fmt.Println(status.LastGameStatus)
	// 		gameStillGoing = false
	// 	}

	// 	if status.LastGameStatus == "ended" || status.LastGameStatus == "win" || status.LastGameStatus == "lose" || status.GameStatus == "false" {
	// 		fmt.Println(status.LastGameStatus)
	// 		gameStillGoing = false
	// 	}
	// 	log.Println("i'm here")
	// 	fireCoord := <-ch
	// 	fmt.Println(fireCoord)
	// 	log.Println("i'm here")
	// }
	// log.Fatal(":(")

}

func waitFunction(ctx context.Context, client connection.Client) (connection.GameStatus, error) {
	waitingForResponse := true
	gotStatus := connection.GameStatus{}
	for waitingForResponse {
		status, err := client.GetStatus(ctx)
		if err != nil {
			log.Println(err)
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
