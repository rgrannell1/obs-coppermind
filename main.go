package main

import (
	copper "github.com/rgrannell1/coppermind/pkg"
)

func main() {
	if err := copper.Coppermind(); err != nil {
		panic(err)
	}
}
