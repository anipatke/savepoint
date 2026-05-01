package main

import (
	"github.com/opencode/savepoint/internal/board"
)

func main() {
	if err := board.Run(); err != nil {
		panic(err)
	}
}