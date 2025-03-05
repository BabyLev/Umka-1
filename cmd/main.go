package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/BabyLev/Umka-1/internal/clients/r4uab"
	"github.com/BabyLev/Umka-1/internal/config"
	"github.com/BabyLev/Umka-1/internal/jobs"
	satellitesRepo "github.com/BabyLev/Umka-1/internal/repo/satellites"
	"github.com/BabyLev/Umka-1/internal/router"
	"github.com/BabyLev/Umka-1/internal/service"
	"github.com/BabyLev/Umka-1/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Config load err")
	}

	storage := storage.New()

	conn, err := pgx.Connect(ctx, cfg.PgConnStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	repo := satellitesRepo.New(conn)

	r4uabClient := r4uab.New(cfg.R4uabURL)
	service := service.New(storage, r4uabClient, repo)
	router := router.SetupRouter(service)
	jobs := jobs.New(storage, r4uabClient)

	done := make(chan struct{})
	defer func() {
		done <- struct{}{}
	}()

	go jobs.Start(done)

	fmt.Println("Server running on localhost:8080")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), router)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
