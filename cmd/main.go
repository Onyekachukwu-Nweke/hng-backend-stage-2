package main

import (
	"log"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/handlers"
)

func main() {
	handler := handlers.NewHandler()

	if err := handler.Serve(); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
