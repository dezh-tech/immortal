package main

import "github.com/dezh-tech/immortal/relay"

func main() {
	s := relay.NewRelay()
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
