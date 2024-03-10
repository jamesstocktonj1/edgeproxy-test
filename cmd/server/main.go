package main

import (
	"log"

	"github.com/jamesstocktonj1/edgeproxy-test/app"
)

func main() {
	s := app.Server{}

	err := s.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Stop()
}
