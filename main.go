package main

import "fmt"

func StorePinboardBookmarks() {
	pb, err := NewPinboardClient()
	if err != nil {
		panic(err)
	}

	bookmarks, err := pb.GetBookmarks()

	if err != nil {
		panic(err)
	}

	fmt.Println(bookmarks)
}

func StoreYoutubeLikes() {

}

func StoreYoutubeMusicLikes() {

}

func main() {
	go StorePinboardBookmarks()
	go StoreYoutubeLikes()
	go StoreYoutubeMusicLikes()
}
