package main

import (
	"log"
	"mini-redis/internal/server"
	"mini-redis/internal/store"
)

func main() {
	s := store.New()
	srv := server.New(s)

	if err := srv.Start(":6380"); err != nil {
		log.Fatal(err)
	}
}
