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

	m.peers[1].Perform(FUP(OK, WAIT, "123"))
	m.peers[0].Perform(FUP(OK, GO, "123"))

	return move(0)
}

func endGame(winnerIndex, loserIndex int) gameFn {

	return func(m *Match) gameFn {

		m.peers[loserIndex].Perform(FUP(END, LOSE))
		m.peers[winnerIndex].Perform(FUP(END, WIN))
		return nil
	}
}

func move(playerIndex int) gameFn {
	
	return func(m *Match) gameFn {
    	
    	//c, err := m.Expect(MOVE, playerIndex)

    	p, _ := m.NextCommand()

		if p == m.peers[playerIndex % 2] {

			if c.isCommand(MOVE) {


			}

			p.Perform(FUP("OK"))
			
			if playerIndex == 2 {
				return endGame(playerIndex % 2, (playerIndex + 1) % 2) 
			}

			return move(playerIndex + 1)
		}

		p.Perform(FUP(ERR_NOT_YOUR_TURN))
		return move(playerIndex) 
	}
}