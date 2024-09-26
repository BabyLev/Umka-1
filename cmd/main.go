package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BabyLev/Umka-1/internal/service"
)

func main() {
	router := service.SetupRouter()

	fmt.Println("Server running on localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Server failed: %w", err)
	}
}
