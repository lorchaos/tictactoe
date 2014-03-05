package server

import ("fmt"
		"strings"
		"strconv"
		"log")

const(
	OK string = "OK"
	ERR_INVALID_SEQUENCE = "ERR_INVALID_SEQUENCE"
  	ERR_INVALID_COMMAND  = "ERR_INVALID_COMMAND"
  	ERR_NOT_YOUR_TURN = "ERR_NOT_YOUR_TURN" 
	ERR_INVALID_MOVE  = "ERR_INVALID_MOVE"
)

const(
	JOIN string = "JOIN"
	MOVE  = "MOVE"
	OP_MOVE = "OP_MOVE"
	QUIT  = "QUIT"
	END  = "END"
    BAN = "BAN"
    WAIT = "WAIT"
    GO = "GO"
    WIN = "WIN"
    LOSE = "LOSE"
)

type board [9] *Peer

type gameFn func(m *Match) gameFn

func FUP(key string, param ... string) *Command {

	p := strings.Join(param, " ")

	c := new(Command)
	c.payload = fmt.Sprintf("%s %s\n", key, p)
	return c
}



func Checker(m *Match) {

	for f := start; f != nil; {
		f = f(m)
	}

	m.Close()
}


func start(m *Match) gameFn {

	waiting := m.peers[1]
	playing := m.peers[0]

	waiting.Perform(FUP(OK, WAIT, "123"))
	playing.Perform(FUP(OK, GO, "123"))

	return move(playing, waiting, new(board))
}

func endGame(winner, loser *Peer) gameFn {

	return func(m *Match) gameFn {

		loser.Perform(FUP(END, LOSE))
		winner.Perform(FUP(END, WIN))
		return nil
	}
}

func process(b *board, move int, p *Peer) bool {

	if move > 0 && move < len(b) && b[move] == nil {
		b[move] = p
		return true
	} 

	return false
}

func done(b *board) bool {

	for i := 0; i < len(b); i++ {

		fmt.Printf("%d : %v\n", i, b[i])
	}

	//TODO calculate here if the board is done
	return false
}


func move(playing, waiting *Peer, b *board) gameFn {
	
	return func(m *Match) gameFn {

    	p, c := m.NextCommand()

    	fmt.Printf("Got %v\n", c)

    	if p == playing {

    		if m, err := strconv.Atoi(c.params[0]); err == nil {

    			if process(b, m, p) {
 
 	   				playing.Perform(FUP(OK))
    				waiting.Perform(FUP(OP_MOVE, strconv.Itoa(m)))

    				if done(b) {
    					return endGame(playing, waiting); 
    				} 
    				return move(waiting, playing, b)
    		
    			} else {
    				p.Perform(FUP(ERR_INVALID_MOVE))
    			}
    		}

		} else {
			p.Perform(FUP(ERR_NOT_YOUR_TURN))
		}

		return move(playing, waiting, b)
	}
}