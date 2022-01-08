package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type PinboardClient struct {
	key string
}

type Bookmark struct {
	Description string `json:"description"`
	Extended    string `json:"extended"`
	Hash        string `json:"hash"`
	Href        string `json:"href"`
	Meta        string `json:"meta"`
	Shared      string `json:"shared"`
	Tags        string `json:"tags"`
	Time        string `json:"time"`
	Toread      string `json:"toread"`
}

type PinboardResponse []Bookmark

/*
 * Construct a Pinboard client
 *
 */
func NewPinboardClient() (*PinboardClient, error) {
	key := os.Getenv("PINBOARD_API_KEY")

	if len(key) == 0 {
		return nil, errors.New("pinboard key missing")
	}

	return &PinboardClient{key}, nil
}

type Result [K any]struct {
	Value K
	Error error
}

/*
 * Fetch all bookmarks from Pinboard
 *
 */
func (pin *PinboardClient) GetBookmarks() chan Result[*Bookmark] {
	result := make(chan Result[*Bookmark])

	go func(){
		start := 0
		offset := PINBOARD_OFFSET_SIZE
		url := "https://api.pinboard.in/v1/posts/all?start=" + fmt.Sprint(start) + "&results=" + fmt.Sprint(offset) + "&format=json&auth_token=" + pin.key

		var bookmarks PinboardResponse

		for {
			res, err := http.Get(url)
			if err != nil {
				result <- Result[*Bookmark]{
					Value: nil,
					Error: err,
				}

				return
			}

			start += offset
			body, err := ioutil.ReadAll(res.Body)

			if res.StatusCode == 429 {
				result <- Result[*Bookmark]{
					Value: nil,
					Error: errors.New("too_many_requests: " + string(body)),
				}

				return

			} else if res.StatusCode != 200 {
				result <- Result[*Bookmark]{
					Value: nil,
					Error: errors.New("bad_response: " + string(body)),
				}

				return
			}

			if err != nil {
				result <- Result[*Bookmark]{
					Value: nil,
					Error: err,
				}

				return
			}

			var data PinboardResponse
			err = json.Unmarshal(body, &data)
			if err != nil {
				result <- Result[*Bookmark]{
					Value: nil,
					Error: err,
				}

				return
			}

			if len(data) == 0 {
				break
			}

			for _, bookmark := range bookmarks {
				fmt.Println(bookmark)
				result <- Result[*Bookmark]{
					Value: &bookmark,
					Error: nil,
				}
			}
		}
	}()

	return result
}
