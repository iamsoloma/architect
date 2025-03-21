package main

import (
	"architect/api"
	"log"
)

func main() {
	config := api.Config{
		ListenAddr: "0.0.0.0:1304",
	}

	server, err := api.NewServer(config)
	if err!=nil{
		log.Fatal(err)
	}

	
	server.Start()
}