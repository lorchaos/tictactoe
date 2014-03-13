package server

import (
	"log"
	"net"
	"github.com/lorchaos/tictactoe/peer"
)

type Server struct {
	listener net.Listener
}


func NewServer() *Server {
	return new(Server)
}

func (server *Server) Stop() {

	server.listener.Close()
}

func (server *Server) Start() chan *peer.Peer {

	ln, err := net.Listen("tcp", ":2020")
	if err != nil {
		// handle error
	}

	server.listener = ln

	c := make(chan *peer.Peer)
	
	go func() {

		for {

			conn, err := ln.Accept()

			if err != nil {

				log.Printf("Unable to accept new connections: %v", err)
				close(c)
				return

			} else {

				c <- peer.NewPeer(conn)
			}
		}
	}()

	return c
}

