package server

import ("fmt"
		"strings")

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

	return move(playing, waiting)
}

func endGame(winner, loser *Peer) gameFn {

	return func(m *Match) gameFn {

		loser.Perform(FUP(END, LOSE))
		winner.Perform(FUP(END, WIN))
		return nil
	}
}

func process() (bool, bool) {
	return false, true
}


func move(playing, waiting *Peer) gameFn {
	
	return func(m *Match) gameFn {

    	p, c := m.NextCommand()

    	fmt.Printf("Got %v\n", c)

    	if p == playing {

	    	done, valid := process()

    		if valid {
 
 	   			playing.Perform(FUP(OK))
    			waiting.Perform(FUP(OP_MOVE, "0"))

    			if done {
    				return endGame(playing, waiting); 
    			}
    			return move(waiting, playing)
    		
    		} else {
    			p.Perform(FUP(ERR_INVALID_MOVE))
    		}

		} else {
			p.Perform(FUP(ERR_NOT_YOUR_TURN))
		}

		return move(playing, waiting)
	}
}