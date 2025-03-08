package main

import (
	"architect/api"
	"log"
)

func main() {
	config := api.Config{
		ListenAddr: "0.0.0.0:1304",
		AddrOfContainerDaemonSocket: "/run/containerd/containerd.sock",
	}

	server, err := api.NewServer(config)
	if err!=nil{
		log.Fatal(err)
	}

	
	server.Start()
}