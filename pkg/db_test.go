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
