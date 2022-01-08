package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func StorePinboardBookmarks(db *CoppermindDb) error {
	pb, err := NewPinboardClient()
	if err != nil {
		panic(err)
	}

	bookmarkResults := pb.GetBookmarks()

	if err != nil {
		return err
	}

	tx, _ := db.Db.Begin()

	for result := range bookmarkResults {
		fmt.Println(result)
		err := result.Error
		bookmark := result.Value

		if err != nil {
			return err
		}

		err = db.InsertBookmark(tx, bookmark)

		if err != nil {
			return err
		}
	}

	return nil
}

func Coppermind() error {
	home, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	dbpath := filepath.Join(home, ".diatom.sqlite")
	db, err := NewCoppermindDb(dbpath)
	if err != nil {
		panic(err)
	}

	defer db.Db.Close()

	err = db.CreateTables()
	if err != nil {
		panic(err)
	}

	err = StorePinboardBookmarks(db)
	if err != nil {
		panic(err)
	}

	return nil
}
