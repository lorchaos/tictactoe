package server

import (
	"time"
	"github.com/lorchaos/tictactoe/peer"
)

func MatchBuilder(c chan *peer.Peer) chan *peer.Match {

	output := make(chan *peer.Match)
	
	go func() {

		id := 0
		match := newMatch(id)

		for {

			select {
			case p := <- c:
				match.AddPeer(p)

				if match.IsComplete() {
					output <- match
					id = id + 1
					match = newMatch(id)
				}

			case <-time.After(5 * time.Second):
				match.Broadcast(peer.NewCommand("Waiting for peer"))

			}
		}
	}()

	return output
}

func newMatch(id int) *peer.Match {

	m := new(peer.Match)
	m.Peers = make([]*peer.Peer, 0, 2)
	m.Id = id
	return m
}


