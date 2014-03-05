package protocol

const(
	OK string = "OK"
	ERR_INVALID_SEQUENCE = "2"
  	ERR_INVALID_COMMAND  = "4"
  	ERR_NOT_YOUR_TURN = "8" 
	ERR_INVALID_MOVE  = "16"
)

const(
	JOIN string = "JOIN"
	MOVE  = "MOVE"
	OP_MOVE = "OP_MOVE"
	QUIT  = "QUIT"
	END  = "END"
    BAN = "BAN"
)