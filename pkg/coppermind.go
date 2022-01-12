package coppermind

import (
	"os"
	"path/filepath"
)

func StorePinboardBookmarks(db *CoppermindDb) error {
	pb, err := NewPinboardClient()
	if err != nil {
		return err
	}

	lastUpdate, err := pb.GetLastUpdate()
	if err != nil {
		return err
	}

	changed, err := db.PinboardChanged(lastUpdate)
	if err != nil {
		return err
	}

	if !changed {
		return nil
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

	return db.UpdateLastUpdated(lastUpdate)
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
