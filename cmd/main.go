package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/BabyLev/Umka-1/internal/clients/r4uab"
	"github.com/BabyLev/Umka-1/internal/config"
	"github.com/BabyLev/Umka-1/internal/jobs"
	locationsRepo "github.com/BabyLev/Umka-1/internal/repo/locations"
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

	repoSats := satellitesRepo.New(conn)
	repoLocs := locationsRepo.New(conn)

	r4uabClient := r4uab.New(cfg.R4uabURL)
	service := service.New(r4uabClient, repoSats, repoLocs)
	router := router.SetupRouter(service)

	// TODO: repo
	jobs := jobs.New(storage, r4uabClient, repoSats)
	go jobs.Start(ctx)

	fmt.Println("Server running on localhost:8080")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), router)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}
