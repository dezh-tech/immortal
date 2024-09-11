package main

import "github.com/dezh-tech/immortal/relay"

// TODO::: create a full functioning CLI to manage rely.

func main() {
	s := relay.NewRelay()
	err := s.Start()
	if err != nil {
		panic(err)
	}
}
