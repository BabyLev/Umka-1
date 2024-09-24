package satellite

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/joshuaferrara/go-satellite"
)

type Satellite struct {
	line1 string
	line2 string
	name  string

	sat *satellite.Satellite
}

type SatelliteCoords struct {
	Lat float64
	Lon float64
	Alt float64
}

func New(line1 string, line2 string, name string) Satellite {
	sat := satellite.TLEToSat(line1, line2, satellite.GravityWGS84)

	return Satellite{
		line1: line1,
		line2: line2,
		name:  name,
		sat:   &sat,
	}
}

// функция возвращает широту, долготу, высоту спутника
func (s Satellite) Calculate(t time.Time) (*SatelliteCoords, error) {
	if s.sat == nil {
		return nil, errors.New("sateliite is not configured")
	}

	// рассчитываем позицию спутника на переданный момент времени
	position, _ := satellite.Propagate(*s.sat, t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	// GST
	// вернет значение времени в радианах, угловое положение Земли относительно полярной звезды на основе переданного времени
	gst := satellite.GSTimeFromDate(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	// Geodesic coordinates (ECI -> LLA)
	alt, _, latLng := satellite.ECIToLLA(position, gst)

	// randians to degress
	lat := latLng.Latitude * satellite.RAD2DEG
	lon := latLng.Longitude * satellite.RAD2DEG

	// долгота будет в диапазоне между -180 и +180
	lon = math.Mod(lon+540, 360) - 180

	lat = math.Mod(lat+540, 360) - 180

	return &SatelliteCoords{
		Lat: lat,
		Lon: lon,
		Alt: alt,
	}, nil
}

func (s *Satellite) UpdateTLE(line1, line2 string) {
	s.line1 = line1
	s.line2 = line2

	sat := satellite.TLEToSat(line1, line2, satellite.GravityWGS84)
	s.sat = &sat
}

func (sc SatelliteCoords) LatDirection() string {
	if sc.Lat >= 0 {
		return "N"
	}

	return "S"
}

func (sc SatelliteCoords) LonDirection() string {
	if sc.Lon >= 0 {
		return "E"
	}

	return "W"
}

func (sc SatelliteCoords) String() string {
	return fmt.Sprintf("%.6f°%s %.6f°%s", math.Abs(sc.Lat), sc.LatDirection(), math.Abs(sc.Lon), sc.LonDirection())
}
