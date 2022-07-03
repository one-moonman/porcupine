package main

import (
	"bug-free-octo-broccoli/api"
	"bug-free-octo-broccoli/storage"
)

func main() {
	defer storage.CloseJsonStorage()
	defer storage.CloseMemoryStorage()
	server := new(api.Server)
	server.Run()
}
