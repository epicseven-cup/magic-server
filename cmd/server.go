package main

import (
	"github.com/epicseven-cup/magic-server/pkg"
	"log"
)

func main() {
	server := pkg.NewServer()
	err := server.Run()
	if err != nil {
		log.Fatalln(err)
		return
	}
}
