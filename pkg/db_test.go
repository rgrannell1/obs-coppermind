package coppermind

import "testing"

func TestNewCoppermindDb(test *testing.T) {
	_, err := NewCoppermindDb(":memory:")
	if err != nil {
		test.Errorf("error creating coppermind db: %v", err)
	}
}

func TestCreateTables(test *testing.T) {
	db, err := NewCoppermindDb(":memory:")
	if err != nil {
		test.Errorf("error creating coppermind db: %v", err)
	}

	err = db.CreateTables()
	if err != nil {
		test.Errorf("error creating tables: %v", err)
	}

	// check tables exist
}

func TestDropBookmarks(test *testing.T) {
	db, err := NewCoppermindDb(":memory:")
	if err != nil {
		test.Errorf("error creating coppermind db: %v", err)
	}

	// create tables first
	err = db.CreateTables()
	if err != nil {
		test.Errorf("error creating tables: %v", err)
	}

	// drop

	err = db.DropBookmarks()
	if err != nil {
		test.Errorf("error dropping bookmarks: %v", err)
	}
}

func TestInsertBookmark(test *testing.T) {

}

func TestUpdateLastUpdated(test *testing.T) {

}

func TestGetLastUpdated(test *testing.T) {

}

func TestPinboardChanged(test *testing.T) {

}
