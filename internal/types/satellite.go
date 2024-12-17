package types

import "github.com/BabyLev/Umka-1/satellite"

type Satellite struct {
	satellite.Satellite

	Name    string
	NoradID *int
}

// координаты наблюдателя
type ObserverLocation struct {
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

type Location struct {
	Lon float64 `json:"lon"` // долгота
	Lat float64 `json:"lat"` // широта
	Alt float64 `json:"alt"` // высота
}
