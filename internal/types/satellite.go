package types

import "github.com/BabyLev/Umka-1/satellite"

type Satellite struct {
	satellite.Satellite

	Name    string
	NoradID *int
}
