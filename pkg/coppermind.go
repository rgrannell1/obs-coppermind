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

	// retreive pinboard last-update time
	lastUpdate, err := pb.GetLastUpdate()
	if err != nil {
		return err
	}

	// check whether pinboard has new bookmarks
	changed, err := db.PinboardChanged(lastUpdate)
	if err != nil {
		return err
	}

	// bookmarks are the same as last run, so
	// no updates are required
	if !changed {
		return nil
	}

	// drop existing bookmarks
	if err = db.DropBookmarks(); err != nil {
		return err
	}

	// enumerate through bookmarks from pinboard, and
	// insert each into the database
	bookmarkResults := pb.GetBookmarks()

	for result := range bookmarkResults {
		if err := result.Error; err != nil {
			return err
		}

		if err = db.InsertBookmark(result.Value); err != nil {
			return err
		}
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

	err = StoreGithubStars(db)
	if err != nil {
		panic(err)
	}

	return nil
}
