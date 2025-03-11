package main

import (
	"log"

	"github.com/oskov/cqlshvm/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
