package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BabyLev/Umka-1/internal/router"
	"github.com/BabyLev/Umka-1/internal/service"
	"github.com/BabyLev/Umka-1/internal/storage"
)

func main() {
	storage := storage.New()
	service := service.New(storage)
	router := router.SetupRouter(service)

	fmt.Println("Server running on localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Server failed: %w", err)
	}
}
