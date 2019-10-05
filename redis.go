package main

import (
	"fmt"
	"log"
	"strconv"

	// Import the redigo/redis package.
	"github.com/gomodule/redigo/redis"
)

// Define a custom struct to hold Album data.
type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func main() {
	// Establish a connection to the Redis server listening on port
	// 6379 of the local machine. 6379 is the default port, so unless
	// you've already changed the Redis configuration file this should
	// work.
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("===========")
	// Importantly, use defer to ensure the connection is always
	// properly closed before exiting the main() function.
	defer conn.Close()

	reply, err := redis.StringMap(conn.Do("HGETALL", "album:1"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("===========")

	// Use the populateAlbum helper function to create a new Album
	// object from the map[string]string.
	fmt.Println(reply)
	album, err := populateAlbum(reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", album)
}

func populateAlbum(reply map[string]string) (*Album, error) {
	var err error
	album := new(Album)
	album.Title = reply["title"]
	album.Artist = reply["artist"]
	// We need to use the strconv package to convert the 'price' value
	// from a string to a float64 before assigning it.
	album.Price, err = strconv.ParseFloat(reply["price"], 64)
	if err != nil {
		return nil, err
	}
	// Similarly, we need to convert the 'likes' value from a string to
	// an integer.
	album.Likes, err = strconv.Atoi(reply["likes"])
	if err != nil {
		return nil, err
	}
	return album, nil
}
