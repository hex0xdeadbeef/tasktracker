package main

import (
	"log"
	"tasktracker/internal/app"
)

// main starts an application calling app.Run()
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
