package server

import ("time")

type Match struct {

	peers []*Peer
}

type GameRunner func(m *Match)

type matchMaker struct {
	players chan *Peer
}

func (m *matchMaker) AddPlayer(p *Peer) {
	m.players <- p
}


func (m *matchMaker) loop(g GameRunner) {

	match := NewMatch()

	for {

		select {
			case p := <- m.players:
				match.AddPeer(p)

				if match.IsComplete() {

					go g(match)
					
					match = NewMatch()
				}

				case <- time.After(5 * time.Second):
					match.Broadcast(FUP("Waiting for peer"))

		}
	}
}

func RunMatchMaker(g GameRunner) *matchMaker {

	m := new(matchMaker)
	m.players = make(chan *Peer)

	go m.loop(g)

	return m
}


func NewMatch() *Match {

	m := new(Match)
	m.peers = make([]*Peer, 0, 2)
	return m
}

// do we have all peers?
func (m Match) IsComplete() bool {

	return len(m.peers) == 2
}


func (m *Match) AddPeer(p *Peer) bool {

	if m.IsComplete() {
		return false
	}

	m.peers = append(m.peers, p)

	return true
}

func(m *Match) Expect(c string, p int) (*Peer, Command) { 

	//TODO fix this
	select {

	case c := <-m.peers[0].out:
		return m.peers[0], c

	case c := <-m.peers[1].out:
		return m.peers[1], c

	};
}


func(m *Match) NextCommand() (*Peer, Command) { 

	//TODO fix this
	select {

	case c := <-m.peers[0].out:
		return m.peers[0], c

	case c := <-m.peers[1].out:
		return m.peers[1], c

	};
}


func (m *Match) Broadcast(c *Command) {

	for _, p := range m.peers {

		p.Perform(c)
	}
}