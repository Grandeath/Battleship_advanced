// menu create main menu of the game
package menu

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	setships "github.com/Grandeath/Battleship_advanced/application/setShips"
	"github.com/Grandeath/Battleship_advanced/connection"
)

// UserIntent represents the user's intention from the main menu
type UserIntent uint8

const (
	// WaitForChallenge indicates the user wants to wait for a challenge
	WaitForChallenge UserIntent = iota
	// ExitTheGame indicates the user wants to exit the game
	ExitTheGame
	// StartGame indicates the user wants to start a game
	StartGame
)

// MainMenu displays the main menu and handles user input
func MainMenu(ctx context.Context, client connection.Client) UserIntent {
	var startingHeader connection.StartingHeader

	for {
		fmt.Println("1. Choose a nick")
		fmt.Println("2. Position ship by yourself")
		fmt.Println("3. Show ladderboard")
		fmt.Println("4. Play against bot")
		fmt.Println("5. Play against chosen player")
		fmt.Println("6. Wait for challenge")
		fmt.Println("7. Game manual")
		fmt.Println("8. Exit game")
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Chose a number")

		var question string
		var chosenMenu int

		// Read user input until a valid menu option is chosen
		for {
			if scanner.Scan() {
				question = scanner.Text()
			} else {
				log.Println(scanner.Err())
			}
			var err error
			chosenMenu, err = strconv.Atoi(question)
			if err != nil {
				log.Println(err)
			} else if chosenMenu < 0 || chosenMenu > 8 {
				log.Println("Wrong number")
			} else {
				break
			}
		}
		// Handle the chosen menu option
		switch chosenMenu {
		case 1:
			err := choseNick(scanner, &startingHeader)
			if err != nil {
				log.Println("Not a number")
			}
		case 2:
			setships.StartPositionBoard(ctx, &startingHeader)
		case 3:
			err := ShowLeaderBoard(ctx, client, scanner, startingHeader)
			if err != nil {
				log.Println(err)
			}
		case 4:
			playAgainstBot(&startingHeader)
			client.SetStartingHeader(startingHeader)
			return StartGame
		case 5:
			err := choosePlayer(ctx, client, scanner, &startingHeader)
			if err != nil {
				log.Println(err)
			} else if len(startingHeader.TargetNick) > 0 {
				startingHeader.Wpbot = false
				client.SetStartingHeader(startingHeader)
				return StartGame
			}
		case 6:
			waitForChallenge(&startingHeader)
			client.SetStartingHeader(startingHeader)
			return WaitForChallenge
		case 7:
			err := PrintManual(scanner)
			if err != nil {
				log.Println("Not a number")
			}
		case 8:
			return ExitTheGame
		}
		CallClear()
	}
}

// choseNick prompts the user to choose a nick and description
func choseNick(scanner *bufio.Scanner, startingHeader *connection.StartingHeader) error {
	fmt.Println("Write your Nick")

	var yourDescription string
	var yourNick string

	// Read user input for nick until valid length is entered
	for {
		if scanner.Scan() {
			yourNick = scanner.Text()
		} else {
			return scanner.Err()
		}
		if len(yourNick) > 20 {
			fmt.Println("Nick too long")
			continue
		} else if len(yourNick) == 0 {
			fmt.Println("Nick too short")
			continue
		} else {
			startingHeader.Nick = yourNick
			break
		}
	}
	// Read user input for description until valid length is entered
	for {
		fmt.Println("Write your Desc")

		if scanner.Scan() {
			yourDescription = scanner.Text()
		} else {
			return scanner.Err()
		}
		if len(yourDescription) > 120 {
			fmt.Println("Description too long")
			continue
		} else if len(yourNick) == 0 {
			fmt.Println("Description too short")
			continue
		} else {
			startingHeader.Desc = yourDescription
			break
		}
	}
	return nil
}

// choosePlayer prompts the user to choose a player to challenge
func choosePlayer(ctx context.Context, client connection.Client, scanner *bufio.Scanner, startingHeader *connection.StartingHeader) error {
	playerList, err := client.GetPlayerList(ctx)
	if err != nil {
		return err
	}
	if len(playerList) == 0 {
		fmt.Println("No enemies to challenge returning to menu")
		time.Sleep(time.Millisecond * 1000)
		return nil
	}
	for index, player := range playerList {
		fmt.Printf("%d %s - %s \n", index, player.Nick, player.GameStatus)
	}
	fmt.Println("Choose a player number(-1 to wait for opponent)")
	var chosenPlayer string
	if scanner.Scan() {
		chosenPlayer = scanner.Text()
	} else {
		return scanner.Err()
	}
	numPlayer, err := strconv.Atoi(chosenPlayer)
	if numPlayer == -1 {
		return nil
	}
	if err != nil {
		return err
	}
	if numPlayer >= len(playerList) || numPlayer < 0 {
		return errors.New("num outside the list")
	}

	startingHeader.TargetNick = playerList[numPlayer].Nick

	return nil
}

// playAgainstBot sets the game mode to play against a bot
func playAgainstBot(startingHeader *connection.StartingHeader) {
	startingHeader.Wpbot = true
}

// waitForChallenge sets the game mode to wait for a challenge
func waitForChallenge(startingHeader *connection.StartingHeader) {
	startingHeader.Wpbot = false
	startingHeader.TargetNick = ""
}

// ShowLeaderboard displays the leaderboard with player statistics
func ShowLeaderBoard(ctx context.Context, client connection.Client, scanner *bufio.Scanner, startingHeader connection.StartingHeader) error {
	statsList, err := client.GetLeaderBoard(ctx)
	if err != nil {
		return err
	}

	var playerStats connection.StatsPlayer

	if len(startingHeader.Nick) != 0 {
		playerStats, err = client.GetPlayerScore(ctx, startingHeader.Nick)
		if err != nil {
			return err
		}
	}

	for _, stat := range statsList.GotStats {
		fmt.Println(stat.Rank)
		fmt.Printf("Nick : %s\n", stat.Nick)
		fmt.Printf("Games : %d\n", stat.Games)
		fmt.Printf("Points : %d\n", stat.Points)
		fmt.Printf("Wins : %d\n", stat.Wins)
		fmt.Println()
	}

	if playerStats.GotStat.Rank > 10 {
		fmt.Println(playerStats.GotStat.Rank)
		fmt.Printf("Nick : %s\n", playerStats.GotStat.Nick)
		fmt.Printf("Games : %d\n", playerStats.GotStat.Games)
		fmt.Printf("Points : %d\n", playerStats.GotStat.Points)
		fmt.Printf("Wins : %d\n", playerStats.GotStat.Wins)
		fmt.Println()
	}

	fmt.Println("Press anything to return to menu")
	if scanner.Scan() {
		scanner.Text()
	} else {
		return scanner.Err()
	}

	return nil
}

// PrintManual print game manual
func PrintManual(scanner *bufio.Scanner) error {
	fmt.Println("How to play")
	fmt.Println("1. Choose a nick to save your progress write Your nick in option 1 of menu otherwise You will get random nick.")
	fmt.Println("2. Position ship by yourself - click to position your ships for the next battle")
	fmt.Println("3. Show ladderboard- showing leaderboard of best players. If you specify your nick you can see your results")
	fmt.Println("4. Play against bot - challenge a bot")
	fmt.Println("5. Play against chosen player - challenge a player")
	fmt.Println("6. Wait for challenge - wait for someone to challengu you")
	fmt.Println("7. Game manual - see manual")
	fmt.Println("8. Exit game - quit the game")
	fmt.Println()
	fmt.Println()
	fmt.Println("How to play")
	fmt.Println("To play this game you need to click on the oppenent board to make a shot")
	fmt.Println("When you hit correctly board will show H symbol and let you shot again")
	fmt.Println("You have 60 sec time limit to make a shot otherwise you will lose")
	fmt.Println("This game will be won by person who first will shot all enemy ships")
	fmt.Println()
	fmt.Println()

	fmt.Println("Press anything to return to menu")
	if scanner.Scan() {
		scanner.Text()
	} else {
		return scanner.Err()
	}
	return nil
}
