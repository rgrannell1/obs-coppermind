package main

import (
	"database/sql"
	"fmt"
)

type CoppermindDb struct {
	Db *sql.DB
}

func NewCoppermindDb(fpath string) (*CoppermindDb, error) {
	db, err := sql.Open("sqlite3", fpath)
	if err != nil {
		return &CoppermindDb{}, err
	}

	return &CoppermindDb{db}, nil
}

func (db *CoppermindDb) CreateTables() error {
	// create a file table
	tx, err := db.Db.Begin()
	defer tx.Rollback()

	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS pinboard_bookmark (
		description string,
		extended    string,
		hash        string,
		href        string,
		meta        string,
		shared      string,
		tags        string,
		time        string,
		toread      string
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

func (db *CoppermindDb) InsertBookmark(tx *sql.Tx, bookmark *Bookmark) error {
	fmt.Println(bookmark)
	_, err := tx.Exec(`
	INSERT OR IGNORE INTO pinboard_bookmarks (description, extended, hash, href, meta, shared, tags, time, toread) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, bookmark.Description, bookmark.Extended, bookmark.Hash, bookmark.Href, bookmark.Meta, bookmark.Shared, bookmark.Tags, bookmark.Time, bookmark.Toread)

	if err != nil {
		return err
	}

	return tx.Commit()
}
