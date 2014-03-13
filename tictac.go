package main

import (
	"github.com/lorchaos/tictactoe/server"
	"github.com/lorchaos/tictactoe/game"
	"log"
)

func main() {

	s := server.NewServer()
	c := s.Start()

	for m := range server.MatchBuilder(c) {

		go m.Run(game.Start)

		log.Printf("Match started: %d", m.Id)
	}
}
