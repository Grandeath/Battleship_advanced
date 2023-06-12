// This application handle logic and connection to API server.
// In connection package there is handling of different endpoint connections and create containing struct responses.
// In application pacakage there is a logic, GUI and Ship positioning of the application.
// Whole game is the classic battleship game played with other users or bots on the server.
// Server receive information like playrs nick, shots, ship positioning, person to challenge
// and give responses like shot result, enemy shots, leaderboard and results of the gane.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Grandeath/Battleship_advanced/application"
)

// main create context and pass it to starting app. When game ends function ask if player want to play again.
func main() {
	game := true
	for game {
		ctx := context.Background()
		application.StartApp(ctx)

		scanner := bufio.NewScanner(os.Stdin)
		var question string
		fmt.Println("Do you want to play again? (yes/no)")
		if scanner.Scan() {
			question = scanner.Text()
		} else {
			log.Println(scanner.Err())
		}
		if question == "yes" {
			continue
		} else {
			game = false
		}
	}

}
