package main

import (
	"fmt"
	"os"

	"github.com/opencode/savepoint/internal/board"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}
	if err := board.Run(); err != nil {
		panic(err)
	}
}