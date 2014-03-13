package peer

import ("log"
	"net"
	"fmt"
	"bufio"
"strings")


type MatchRunner func(m *Match) MatchRunner

func (m *Match) Run(runner MatchRunner) {

  for f := runner; f != nil; {
	f = f(m)
  }

  m.Broadcast(NewCommand("BYE"))
}


type Match struct {
	Peers []*Peer
	Id int
}

// do we have all peers?
func (m *Match) IsComplete() bool {

	return len(m.Peers) == 2
}

func (m *Match) AddPeer(p *Peer) bool {

	if m.IsComplete() {
		return false
	}

	m.Peers = append(m.Peers, p)

	return true
}


func (m *Match) Expect(c string, p int) (*Peer, Command) {

	//TODO fix this
	select {

	case c := <-m.Peers[0].out:
		return m.Peers[0], c

	case c := <-m.Peers[1].out:
		return m.Peers[1], c

	}
}

func (m *Match) NextCommand(c string) (*Peer, Command) {

	//TODO fix this
	select {

	case c := <-m.Peers[0].out:
		return m.Peers[0], c

	case c := <-m.Peers[1].out:
		return m.Peers[1], c

	}
}

func (m *Match) Broadcast(c *Command) {

	for _, p := range m.Peers {

		p.Perform(c)
	}
}


type Peer struct {

	// the network connection
	conn net.Conn

	// this is how the peer sends messages to the server
	out chan Command

	// this is how the server communicates with the Peer
	in chan Command
}

type Sendable interface {
	Payload() string
}

type Command struct {
	Sendable
	Id      string
	Params  []string
}

func NewCommand(id string, param ...string) *Command {
	
	c := new(Command)
	c.Id = id
	c.Params = param
	return c
}

func (c *Command) Payload() string  {
	p := strings.Join(c.Params, " ")
	return fmt.Sprintf("%s %s\n", c.Id, p)
}



func (p Peer) Perform(c *Command) {

	p.in <- *c
}

func (p Peer) quit() {

	p.conn.Close()

	close(p.out)
	close(p.in)

	log.Printf("Done.")
}

func NewPeer(conn net.Conn) *Peer {

	//defer conn.Close()

	//TODO create a host abstraction here

	p := new(Peer)
	p.conn = conn
	p.out = make(chan Command)
	p.in = make(chan Command)

	go p.handleWrite()
	go p.handleRead()

	log.Printf("Connected! %v\n", conn)

	return p
}

func (p Peer) handleWrite() {

	w := bufio.NewWriter(p.conn)

	for c := range p.in {

		if c.Id == "BYE" {
			p.quit()
		} else {
			pl := c.Payload()
			fmt.Fprintf(w, pl)
			w.Flush()
		}
	}
}

func (p Peer) handleRead() {

	r := bufio.NewReader(p.conn)
	scanner := bufio.NewScanner(r)

	var c *Command

	for scanner.Scan() {

		t := scanner.Text()

		log.Printf("Handled %s %d\n", t, len(t))

		tokens := strings.Split(t, " ")

		c = new(Command)
		c.Id = tokens[0]
		c.Params = tokens[1:]
		p.out <- *c
	}

	if err := scanner.Err(); err != nil {

		log.Println("reading standard input:", err)

	} else {
		log.Println("no errors")
	}

	//p.quit()

	log.Printf("User %v left.", p.conn)
}
