package main

import (
	"fmt"
	"log"
)

func main() {
	var command string

	for {
		if _, err := fmt.Scan(&command); err != nil {
			panic(err)
		}

		log.Println("command:", command)

		// do command (command)
	}
}
