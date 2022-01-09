package main

import (
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := Coppermind(); err != nil {
		panic(err)
	}
}
