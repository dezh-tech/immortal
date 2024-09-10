package main

import "github.com/dezh-tech/immortal/server"

func main() {
	s := server.NewServer()
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
