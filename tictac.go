package main

import ("github.com/lorchaos/checker/server")

func main() {
	
	server := server.NewServer()

	server.Start()
}