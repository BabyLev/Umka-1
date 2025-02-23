package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BabyLev/Umka-1/internal/clients/r4uab"
	"github.com/BabyLev/Umka-1/internal/jobs"
	"github.com/BabyLev/Umka-1/internal/router"
	"github.com/BabyLev/Umka-1/internal/service"
	"github.com/BabyLev/Umka-1/internal/storage"
)

func main() {
	storage := storage.New()

	r4uabClient := r4uab.New("https://api.r4uab.ru")
	service := service.New(storage, r4uabClient)
	router := router.SetupRouter(service)
	jobs := jobs.New(storage, r4uabClient)

	done := make(chan struct{})
	defer func() {
		done <- struct{}{}
	}()

	go jobs.Start(done)
	jobs.Start(done)

	fmt.Println("Server running on localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Server failed: %w", err)
	}
}
