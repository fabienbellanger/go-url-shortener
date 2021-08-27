package main

import (
	"log"

	"github.com/fabienbellanger/go-url-shortener/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatalln(err)
	}
}
