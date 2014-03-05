package server

import("testing")

func TestCompletion(t *testing.T) {

	m := NewMatch()

	for i := 0; i < 2; i++ {

		if !m.AddPeer(new(Peer)) {
			t.Errorf("Match should not be complete yet")
		}
	}

	if !m.IsComplete() {
		t.Errorf("Match should be complete now")
	}
}

func TestCommand(t *testing.T) {

	t.Error(FUP("OK", "GO", "1234"));
}