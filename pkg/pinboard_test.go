package coppermind

import (
	"os"
	"strings"
	"testing"
)

func TestNewPinboardClient(test *testing.T) {
	currKey := os.Getenv(PINBOARD_KEY_VARIABLE)
	if len(currKey) == 0 {
		test.Errorf("Pinboard key not provided in advance")
	}

	defer os.Setenv(PINBOARD_KEY_VARIABLE, currKey)

	os.Setenv(PINBOARD_KEY_VARIABLE, "")

	if _, err := NewPinboardClient(); err == nil {
		test.Errorf("No error when constructing client, though one was expected: %v", err)
	}

	os.Setenv(PINBOARD_KEY_VARIABLE, currKey)
	client, err := NewPinboardClient()
	if err != nil {
		test.Errorf("Failed %v", err)
	}

	if client.key != currKey {
		test.Errorf("Key did not match %v", err)
	}
}

func TestGetLastUpdate(test *testing.T) {
	currKey := os.Getenv(PINBOARD_KEY_VARIABLE)
	if len(currKey) == 0 {
		test.Errorf("Pinboard key not provided in advance")
	}

	client, err := NewPinboardClient()
	if err != nil {
		test.Errorf("Failed %v", err)
	}

	updateTime, err := client.GetLastUpdate()
	if err != nil {
		test.Errorf("Failed %v", err)
	}

	if len(updateTime) == 0 {
		test.Errorf("Update time not present %v", err)
	}

	if !strings.HasPrefix(updateTime, "202") {
		test.Errorf("GetLastUpdate() did not return a date %v", err)
	}
}

func TestGetBookmarks(test *testing.T) {
	currKey := os.Getenv(PINBOARD_KEY_VARIABLE)
	if len(currKey) == 0 {
		test.Errorf("Pinboard key not provided in advance")
	}

	client, err := NewPinboardClient()
	if err != nil {
		test.Errorf("Failed %v", err)
	}

	count := 0
	for res := range client.GetBookmarks() {
		err := res.Error

		if err != nil {
			test.Errorf("Error retreiving bookmarks %v", err)
		}

		count++
	}

	// -- hard coded for my account
	if count < 1_100 {
		test.Errorf("Too few bookmarks retrieved %v", err)
	}
}
