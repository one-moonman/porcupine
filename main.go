package main

import (
	"bug-free-octo-broccoli/api"
	"bug-free-octo-broccoli/configs"
)

func main() {
	configs.Init()
	server := new(api.Server)
	server.Run()
}
