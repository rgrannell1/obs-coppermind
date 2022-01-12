package coppermind

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

type PinboardBookmarksResponse []Bookmark

type PinboardLastUpdateResponse struct {
	UpdateTime string `json:"update_time"`
}

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

func (pin *PinboardClient) GetLastUpdate() (string, error) {
	url := "https://api.pinboard.in/v1/posts/update?format=json&auth_token="+ pin.key
	res, err := http.Get(url)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 429 {
		err = errors.New("too_many_requests: " + string(body))
	} else if res.StatusCode != 200 {
		err = errors.New("bad_status_code: " + string(body))
	} else {
		err = nil
	}

	var data PinboardLastUpdateResponse
	err = json.Unmarshal(body, &data)

	return data.UpdateTime, err
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

		for {
			url := "https://api.pinboard.in/v1/posts/all?start=" + fmt.Sprint(start) + "&results=" + fmt.Sprint(offset) + "&format=json&auth_token=" + pin.key
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

			var data PinboardBookmarksResponse
			err = json.Unmarshal(body, &data)
			if err != nil {
				result <- Result[*Bookmark]{
					Value: nil,
					Error: err,
				}

				return
			}

			if len(data) == 0 {
				close(result)
				break
			}

			for _, bookmark := range data {
				result <- Result[*Bookmark]{
					Value: &bookmark,
					Error: nil,
				}
			}
		}
	}()

	return result
}
