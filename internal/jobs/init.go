package jobs

import (
	"time"

	"github.com/BabyLev/Umka-1/internal/clients/r4uab"
	"github.com/BabyLev/Umka-1/internal/storage"
)

type Jobs struct {
	storage     *storage.Storage
	r4uabClient *r4uab.Client
}

func New(storage *storage.Storage, r4uabClient *r4uab.Client) *Jobs {
	return &Jobs{
		storage:     storage,
		r4uabClient: r4uabClient,
	}
}

func (j *Jobs) Start(done <-chan struct{}) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return

		case <-ticker.C:
			go j.UpdateSatellitesInfo()
		}
	}
}
