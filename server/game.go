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
    BYE = "BYE"
)

type board [9] *Peer

type gameFn func(m *Match) gameFn

func COMMAND(key string, param ... string) *Command {

	p := strings.Join(param, " ")

	c := new(Command)
	c.id = key
	c.payload = fmt.Sprintf("%s %s\n", key, p)
	return c
}



func Checker(m *Match) {

	for f := start; f != nil; {
		f = f(m)
	}

	m.Broadcast(COMMAND(BYE))
}


func start(m *Match) gameFn {

	waiting := m.peers[1]
	playing := m.peers[0]

	waiting.Perform(COMMAND(OK, WAIT, "123"))
	playing.Perform(COMMAND(OK, GO, "123"))

	return move(playing, waiting, new(board))
}

func endGame(winner, loser *Peer) gameFn {

	return func(m *Match) gameFn {

		loser.Perform(COMMAND(END, LOSE))
		winner.Perform(COMMAND(END, WIN))
		return nil
	}
}

func process(b *board, move int, p *Peer) bool {

	if move >= 0 && move < len(b) && b[move] == nil {
		b[move] = p
		return true
	} 

	return false
}

func done(b *board) bool {

	var ver = make([][]int, 7)
	ver[0] = []int {1, 3, 4}
	ver[1] = []int {3}
	ver[2] = []int {3, 2}
	ver[3] = []int {1}
	ver[6] = []int {1}

	for i, v := range ver {

		for _, v2 := range v {
			
			if 	b[i] != nil &&
			 	b[i] == b[i + v2] && 
				b[i] == b[i + v2 + v2] {

				log.Printf("Done\n");
				return true
			}
		}
	}

	return false
}


func move(playing, waiting *Peer, b *board) gameFn {
	
	return func(m *Match) gameFn {

    	p, c := m.NextCommand()

    	fmt.Printf("Got %v\n", c)

    	if p == playing {

    		if len(c.params) <= 0 {

    			p.Perform(COMMAND(ERR_INVALID_COMMAND))

    		} else if m, err := strconv.Atoi(c.params[0]); err == nil {

    			if process(b, m, p) {
 
 	   				playing.Perform(COMMAND(OK))
    				waiting.Perform(COMMAND(OP_MOVE, strconv.Itoa(m)))

    				if done(b) {
    					return endGame(playing, waiting); 
    				} 
    				return move(waiting, playing, b)
    		
    			} else {
    				p.Perform(COMMAND(ERR_INVALID_MOVE))
    			}
    		}

		} else {
			p.Perform(COMMAND(ERR_NOT_YOUR_TURN))
		}

		return move(playing, waiting, b)
	}
}