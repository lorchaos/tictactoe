package game

import ("testing")


func TestInvalidMoves(t *testing.T) {

	p := "a"
	b := new(board)
	
	if b.process(-1, p) || b.process(9, p) {
		t.Error("Out of bounds movements should not be accepted")
	}

	b.process(0, p) 
	if b.process(0, p) {
		t.Error("The same space cannot be allocated twice")
	}

}

func TestBoard(t *testing.T) {

	b := new(board)

	if b.done() {
		t.Error("Board should not be done when starting")
	}
}

func TestFinishBoard(t *testing.T) {

	p1:= "A"

	b:= new(board)
	
	for i := 0; i < 3; i++ {
		b.process(i, p1)
	}

	if !b.done() {
		t.Error("Board should be done")
	}
}

func TestStaleBoard(t *testing.T) {

	b := new(board)
	
	for i := 0; i < 9; i++ {
		
		if b.stale() {
			t.Error("board should not be stale yet")
		}

		b.process(i, i)
	}
	
	if b.done() {
		t.Error("No one should have won")
	}

	if !b.stale() {
		t.Error("table should be stale")
	}

}
