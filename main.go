package main

import (
	"log"

	"github.com/joho/godotenv"
	copper "github.com/rgrannell1/coppermind/pkg"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := copper.Coppermind(); err != nil {
		panic(err)
	}
}
