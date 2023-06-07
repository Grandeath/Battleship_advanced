package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Grandeath/Battleship_advanced/application"
)

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
