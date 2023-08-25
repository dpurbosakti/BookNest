package main

import (
	"book-nest/app"
	"log"
)

func main() {
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
