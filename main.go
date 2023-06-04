package main

import (
	"context"

	"github.com/Grandeath/Battleship_advanced/application"
)

func main() {
	ctx := context.Background()
	application.StartApp(ctx)
}
