package jobs

import (
	"context"
	"time"

	"github.com/BabyLev/Umka-1/internal/clients/r4uab"
	satellitesRepo "github.com/BabyLev/Umka-1/internal/repo/satellites"
	"github.com/BabyLev/Umka-1/internal/storage"
)

type Jobs struct {
	storage     *storage.Storage
	r4uabClient *r4uab.Client
	repoSats    *satellitesRepo.Repo
}

func New(storage *storage.Storage, r4uabClient *r4uab.Client, repo *satellitesRepo.Repo) *Jobs {
	return &Jobs{
		storage:     storage,
		r4uabClient: r4uabClient,
		repoSats:    repo,
	}
}

func (j *Jobs) Start(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			go j.UpdateSatellitesInfo(ctx)
		}
	}
}
