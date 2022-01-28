package coppermind

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
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

		PRIMARY KEY(href)
	)`)

	if err != nil {
		return errors.Wrap(err, "failed creating pinboard_bookmark")
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS kv_metadata (
		key string,
		value string,

		PRIMARY KEY(key)
	)`)

	if err != nil {
		return errors.Wrap(err, "failed creating kv_metadata")
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS github_star (
		name        string not null,
		description string,
		login       string,
		url         string,
		language    string,
		topics      string,

		PRIMARY KEY(name)
	)`)

	if err != nil {
		return errors.Wrap(err, "failed creating kv_metadata")
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
		return errors.Wrap(err, "failed deleting from pinboard")
	}

	return tx.Commit()
}

func (db *CoppermindDb) InsertBookmark(bookmark Bookmark) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT OR REPLACE INTO pinboard_bookmark (description, extended, hash, href, meta, shared, tags, time, toread) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, bookmark.Description, bookmark.Extended, bookmark.Hash, bookmark.Href, bookmark.Meta, bookmark.Shared, bookmark.Tags, bookmark.Time, bookmark.Toread)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *CoppermindDb) InsertStar(star *StarredRepository) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}

	topics, err := json.Marshal(star.Topics)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT OR REPLACE INTO github_star (name, description, login, url, language, topics) VALUES (?, ?, ?, ?, ?, ?)
	`, star.FullName, star.Description, star.Login, star.Url, star.Language, string(topics))
	if err != nil {
		return err
	}

	return tx.Commit()
}

/*
 *
 */
func (db *CoppermindDb) UpdateLastUpdated(lastUpdate string) error {
	tx, _ := db.Db.Begin()
	defer tx.Rollback()

	_, err := tx.Exec(`
	INSERT OR REPLACE INTO kv_metadata (key, value) VALUES ("lastUpdated", ?)
	`, lastUpdate)

	if err != nil {
		return errors.Wrap(err, "failed updating metadata")
	}

	return tx.Commit()
}

/*
 * Get last updated date for pinboard from database
 */
func (db *CoppermindDb) GetLastUpdated() (string, error) {
	row := db.Db.QueryRow(`SELECT value from kv_metadata WHERE key = "lastUpdated"`)

	var value string

	switch err := row.Scan(&value); err {
	case sql.ErrNoRows:
		return "", nil
	case nil:
		return value, nil
	default:
		return "", errors.Wrap(err, "failed selecting lastUpdated")
	}
}

/*
 * Has pinboard changed, based on a stored update value?
 */
func (db *CoppermindDb) PinboardChanged(lastChanged string) (bool, error) {
	stored, err := db.GetLastUpdated()

	if err != nil {
		return false, err
	}

	return stored != lastChanged, nil
}
