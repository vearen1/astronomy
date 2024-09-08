package main

import (
	"astronomy/astronomy/internal/server"
	"log"
)

func main() {
	server := server.NewMinecraftServer("0.0.0.0", 25565)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
