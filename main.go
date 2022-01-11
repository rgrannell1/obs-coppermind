package main

import (
	_ "github.com/mattn/go-sqlite3"
	copper "github.com/rgrannell1/coppermind/pkg"
)

func main() {
	if err := copper.Coppermind(); err != nil {
		panic(err)
	}
}
