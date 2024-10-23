package satellite

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/joshuaferrara/go-satellite"
)

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
		Lat:       lat,
		Lon:       lon,
		Alt:       alt,
		GMapsLink: fmt.Sprintf("https://www.google.com/maps/place/%f+%f", lat, lon),
	}, nil
}

func (s Satellite) LookAngles(t time.Time, obsCoords ObserverCoords) LookAngles {
	jday := satellite.JDay(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	// рассчитываем позицию спутника на переданный момент времени
	satPosition, _ := satellite.Propagate(*s.sat, t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())

	observerPosition := satellite.LatLong{
		Latitude:  obsCoords.Lat,
		Longitude: obsCoords.Lon,
	}

	lookAngles := satellite.ECIToLookAngles(satPosition, observerPosition, obsCoords.Alt, jday)

	return LookAngles{
		Az:    lookAngles.Az,
		El:    lookAngles.El,
		Range: lookAngles.Rg,
	}
}

// от текущего времени  рассчитает временные диапазоны, когда видно спутник над заданной точкой
// в нужном количестве (от 1 до n диапазонов)
// один диапазон - это время восхода и захода спутника
// n >= 1
func (s Satellite) VisibleTimeRange(t time.Time, obsCoords ObserverCoords, n int) []TimeRange {
	// запомнить время, когда элевация равна нулю (но была отрицательная) - значит, спутник появился над горизонтом
	// запомнить время, когда элевация снова стала равна нулю (но была положительная) - значит, спутник опустился за горизонт
	// повторить Н раз

	timeRangeList := make([]TimeRange, 0, n)
	tCalc := t

	for i := 0; i < n; i++ {
		var timeRange TimeRange

		var prevEl float64

		for {
			lookAngles := s.LookAngles(tCalc, obsCoords)

			if prevEl < 0 && lookAngles.El >= 0 {
				timeRange.From = tCalc // спутник появился над горизонтом
			}

			if prevEl > 0 && lookAngles.El <= 0 {
				timeRange.To = tCalc // зашел за горизонт
				diff := timeRange.To.Sub(timeRange.From)
				timeRange.Difference = diff.String()
				timeRangeList = append(timeRangeList, timeRange) // можем добавлять полученный диапазон в список диапазонов
				break
			}

			prevEl = lookAngles.El

			if t.Add(time.Hour * 24 * 5).Before(tCalc) {
				break
			}
		}
	}

	return timeRangeList
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
