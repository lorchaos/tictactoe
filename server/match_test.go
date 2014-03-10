package server

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

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

func TestArray(t *testing.T) {

	var v [4]*testing.T

	fmt.Printf("G %v", v[3])
}

func estClient(t *testing.T) {

	for i := 0; i < 100; i++ {
		go clien(i)
	}

	time.Sleep(2 * time.Minute)

}

func clien(i int) {

	fmt.Printf("Starting client %d\n", i)

	conn, err := net.Dial("tcp", "localhost:2020")
	defer conn.Close()

	if err != nil {
		// handle error
	}
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	s, _ := bufio.NewReader(conn).ReadString('\n')

	fmt.Printf("From server %s", s)

	writer := bufio.NewWriter(conn)

	for i := 0; i < 100000; i++ {

		writer.WriteString("move 1\n")
		writer.Flush()
	}

	writer.Flush()

}
