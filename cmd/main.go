package main

import (
	"net/http"
	"porcupine/pkg/api"
	"porcupine/pkg/config"
	"porcupine/pkg/storage"
)

func main() {
	config.Init()
	//storage.ConnectJsonStorage(config.MongoOptions)
	storage.ConnectCacheStorage(&config.RedisOptions)

	http.ListenAndServe(":3000", api.Router())
}
