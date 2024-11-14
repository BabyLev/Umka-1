package satellite

import (
	"time"

	"github.com/joshuaferrara/go-satellite"
)

type Satellite struct {
	line1 string
	line2 string

	sat *satellite.Satellite
}

// Getter: Line1
func (s *Satellite) GetLine1() string {
	return s.line1
}

// Getter: Line2
func (s *Satellite) GetLine2() string {
	return s.line2
}

type SatelliteCoords struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Alt       float64 `json:"alt"`
	GMapsLink string  `json:"gmapsLink"`
}

type LookAngles struct {
	Az    float64 `json:"az"`
	El    float64 `json:"el"`
	Range float64 `json:"range"`
}

type ObserverCoords struct {
	Lon float64 // долгота
	Lat float64 // широта
	Alt float64 // высота
}

type TimeRange struct {
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	Difference string    `json:"difference"`
}
