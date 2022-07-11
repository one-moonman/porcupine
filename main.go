package main

import (
	"bug-free-octo-broccoli/api"
)

func main() {
	server := new(api.Server)
	server.Run()
}
