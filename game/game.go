package game

import (
	"fmt"
	//"log"
	"strconv"
	"github.com/lorchaos/tictactoe/peer"
)

const (
	OK                   string = "OK"
	ERR_INVALID_SEQUENCE        = "ERR_INVALID_SEQUENCE"
	ERR_INVALID_COMMAND         = "ERR_INVALID_COMMAND"
	ERR_NOT_YOUR_TURN           = "ERR_NOT_YOUR_TURN"
	ERR_INVALID_MOVE            = "ERR_INVALID_MOVE"
	JOIN           = "JOIN"
	MOVE           = "MOVE"
	OP_MOVE        = "OP_MOVE"
	QUIT           = "QUIT"
	END            = "END"
	BAN            = "BAN"
	WAIT           = "WAIT"
	GO             = "GO"
	WIN            = "WIN"
	LOSE           = "LOSE"
	DRAW           = "DRAW"
	BYE            = "BYE"
)
 

func Start(m *peer.Match) peer.MatchRunner {

	waiting := m.Peers[1]
	playing := m.Peers[0]

	id := strconv.Itoa(m.Id)

	waiting.Perform(peer.NewCommand(OK, WAIT, id))
	playing.Perform(peer.NewCommand(OK, GO, id))

	return move(playing, waiting, new(board))
}

func staleGame() peer.MatchRunner {

	return func(m *peer.Match) peer.MatchRunner {
		m.Broadcast(peer.NewCommand(END, DRAW))
		return nil
	}
}


func endGame(winner, loser *peer.Peer) peer.MatchRunner {

 	loser.Perform(peer.NewCommand(END, LOSE))
	winner.Perform(peer.NewCommand(END, WIN))

	return nil
}


func move(playing, waiting *peer.Peer, b *board) peer.MatchRunner {

	return func(m *peer.Match) peer.MatchRunner {

		p, c := m.NextCommand(MOVE)

		fmt.Printf("Got %v\n", c)

		if p == playing {

			if len(c.Params) <= 0 {

				p.Perform(peer.NewCommand(ERR_INVALID_COMMAND))

			} else if m, err := strconv.Atoi(c.Params[0]); err == nil {

				if b.process(m, p) {

					playing.Perform(peer.NewCommand(OK))
					waiting.Perform(peer.NewCommand(OP_MOVE, strconv.Itoa(m)))
					if b.done() {
						return endGame(playing, waiting)
					} else if b.stale() {
						return staleGame()
					}
					
					return move(waiting, playing, b)

				} else {
					p.Perform(peer.NewCommand(ERR_INVALID_MOVE))
				}
			}

		} else {
			p.Perform(peer.NewCommand(ERR_NOT_YOUR_TURN))
		}

		return move(playing, waiting, b)
	}
}
