package main

import (
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

	err = db.DropBookmarks()
	if err != nil {
		return err
	}

	tx, _ := db.Db.Begin()

	for result := range bookmarkResults {
		if result.Error != nil {
			return err
		}

		if err = db.InsertBookmark(tx, result.Value); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
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
