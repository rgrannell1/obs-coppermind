package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := Coppermind()

	if err != nil {
		panic(err)
	}
}
