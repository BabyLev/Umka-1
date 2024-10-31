package satellite

import (
	"time"

	"github.com/joshuaferrara/go-satellite"
)

type Satellite struct {
	line1 string
	line2 string
	name  string

	sat *satellite.Satellite
}

// Getter: Name
func (s *Satellite) GetName() string {
	return s.name
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
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Alt          float64 `json:"alt"`
	GMapsLink    string  `json:"gmapsLink"`
	AnotherField string  `json:"-"`
}

type LookAngles struct {
	Az    float64 `json:"az"`
	El    float64 `json:"el"`
	Range float64 `json:"range"`
}

type ObserverCoords struct {
	Lon float64
	Lat float64
	Alt float64
}

type TimeRange struct {
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	Difference string    `json:"difference"`
}
