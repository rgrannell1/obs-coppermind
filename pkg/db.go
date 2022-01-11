package coppermind

import (
	"database/sql"
)

type CoppermindDb struct {
	Db *sql.DB
}

func NewCoppermindDb(fpath string) (*CoppermindDb, error) {
	db, err := sql.Open("sqlite3", "file:"+fpath+"?_foreign_keys=true&_busy_timeout=5000&_journal_mode=WAL")

	if err != nil {
		return &CoppermindDb{}, err
	}

	return &CoppermindDb{db}, nil
}

func (db *CoppermindDb) CreateTables() error {
	tx, err := db.Db.Begin()
	defer tx.Rollback()

	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS pinboard_bookmark (
		description string,
		extended    string,
		hash        string NOT NULL,
		href        string NOT NULL,
		meta        string,
		shared      string,
		tags        string,
		time        string,
		toread      string,

		PRIMARY KEY(hash)
	)`)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *CoppermindDb) DropBookmarks() error {
	tx, _ := db.Db.Begin()
	defer tx.Rollback()

	_, err := tx.Exec("DELETE FROM pinboard_bookmark")

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *CoppermindDb) InsertBookmark(tx *sql.Tx, bookmark *Bookmark) error {
	_, err := tx.Exec(`
	INSERT OR REPLACE INTO pinboard_bookmark (description, extended, hash, href, meta, shared, tags, time, toread) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, bookmark.Description, bookmark.Extended, bookmark.Hash, bookmark.Href, bookmark.Meta, bookmark.Shared, bookmark.Tags, bookmark.Time, bookmark.Toread)

	if err != nil {
		return err
	}

	return nil
}
