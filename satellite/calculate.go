package satellite

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/joshuaferrara/go-satellite"
)

// Constants for the visibility calculation
const (
	// Default step for coarse search (finding the interval containing rise/set)
	defaultCoarseSearchStep = 5 * time.Minute
	// Desired precision for rise/set time
	defaultEventTimePrecision = time.Second
	// Maximum duration to search forward for the next pass
	defaultMaxSearchDuration = 7 * 24 * time.Hour // Search up to 7 days ahead
	// Minimum elevation considered "visible" (degrees)
	minVisibleElevation = 0.0
)

func New(line1 string, line2 string) Satellite {
	sat := satellite.TLEToSat(line1, line2, satellite.GravityWGS84)

	return Satellite{
		line1: line1,
		line2: line2,
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
		Latitude:  obsCoords.Lat * satellite.DEG2RAD,
		Longitude: obsCoords.Lon * satellite.DEG2RAD,
	}

	lookAngles := satellite.ECIToLookAngles(satPosition, observerPosition, obsCoords.Alt, jday)

	return LookAngles{
		Az:    lookAngles.Az * satellite.RAD2DEG,
		El:    lookAngles.El * satellite.RAD2DEG,
		Range: lookAngles.Rg,
	}
}

// от текущего времени  рассчитает временные диапазоны, когда видно спутник над заданной точкой
// в нужном количестве (от 1 до n диапазонов)
// один диапазон - это время восхода и захода спутника
// n >= 1
// findNextElevationEvent searches for the next time the satellite's elevation crosses
// the minVisibleElevation threshold. It uses a coarse search followed by a
// bisection method for refinement.
// startTime: Time to start searching from.
// obsCoords: Observer's coordinates.
// findRise: If true, search for elevation crossing from negative to positive (rise).
//
//	If false, search for elevation crossing from positive to negative (set).
//
// coarseStep: Step size for the initial coarse search.
// precision: Desired time precision for the event.
// maxDuration: Maximum time duration to search forward.
func (s Satellite) findNextElevationEvent(startTime time.Time, obsCoords ObserverCoords, findRise bool, coarseStep, precision, maxDuration time.Duration) (time.Time, bool) {
	endTime := startTime.Add(maxDuration)
	currentTime := startTime

	// --- Coarse Search ---
	var intervalStartTime, intervalEndTime time.Time
	foundInterval := false

	// Get initial elevation state
	prevLookAngles := s.LookAngles(currentTime, obsCoords)
	prevEl := prevLookAngles.El

	// Helper function to check the crossing condition
	checkCrossing := func(currentEl float64) bool {
		if findRise {
			// Rise: Was below horizon, now at or above
			return prevEl < minVisibleElevation && currentEl >= minVisibleElevation
		}
		// Set: Was at or above horizon, now below
		return prevEl >= minVisibleElevation && currentEl < minVisibleElevation
	}

	for currentTime.Before(endTime) {
		nextTime := currentTime.Add(coarseStep)
		if nextTime.After(endTime) {
			nextTime = endTime
		}

		currentLookAngles := s.LookAngles(nextTime, obsCoords)
		currentEl := currentLookAngles.El

		if checkCrossing(currentEl) {
			intervalStartTime = currentTime
			intervalEndTime = nextTime
			foundInterval = true
			break // Found the interval containing the event
		}

		prevEl = currentEl
		currentTime = nextTime

		if currentTime.Equal(endTime) {
			break // Reached max search duration in coarse search
		}
	}

	if !foundInterval {
		return time.Time{}, false // Event not found within maxDuration
	}

	// --- Fine Search (Bisection Method) ---
	lowTime := intervalStartTime
	highTime := intervalEndTime

	for highTime.Sub(lowTime) > precision {
		midTime := lowTime.Add(highTime.Sub(lowTime) / 2)
		midEl := s.LookAngles(midTime, obsCoords).El

		// Check if the event is in the first half or second half
		// Note: This logic is slightly different for rise vs set
		var conditionMet bool
		if findRise {
			// For rise, if midEl is still below, the event is in the second half
			conditionMet = midEl < minVisibleElevation
		} else {
			// For set, if midEl is still above, the event is in the second half
			conditionMet = midEl >= minVisibleElevation
		}

		if conditionMet {
			lowTime = midTime
		} else {
			highTime = midTime
		}
	}

	// Return the time at the start of the final refined interval
	// (precision is now met)
	return lowTime, true
}

// VisibleTimeRange calculates the next 'n' time ranges when the satellite is visible
// above the minimum elevation from the observer's location.
// t: The time to start searching from.
// obsCoords: Observer's coordinates (latitude, longitude, altitude).
// n: The desired number of visibility ranges to find (n >= 1).
func (s Satellite) VisibleTimeRange(t time.Time, obsCoords ObserverCoords, n int) []TimeRange {
	if n <= 0 {
		return []TimeRange{}
	}

	timeRangeList := make([]TimeRange, 0, n)
	currentTime := t

	// Use default constants, but allow potential future configuration
	coarseStep := defaultCoarseSearchStep
	precision := defaultEventTimePrecision
	maxDuration := defaultMaxSearchDuration // Max duration for *each* rise/set search

	for len(timeRangeList) < n {
		// 1. Find the next rise time
		riseTime, foundRise := s.findNextElevationEvent(currentTime, obsCoords, true, coarseStep, precision, maxDuration)
		if !foundRise {
			// If no more rises found within the max search duration, stop.
			break
		}

		// 2. Find the next set time *after* the rise time
		// Start searching slightly after rise to avoid finding the same event if precision is limited
		// or if rise/set happen very close together.
		searchSetStartTime := riseTime.Add(precision)
		setTime, foundSet := s.findNextElevationEvent(searchSetStartTime, obsCoords, false, coarseStep, precision, maxDuration)

		if !foundSet {
			// This is less likely if a rise was found, but possible if the pass is
			// extremely short or calculation issues occur near the end of maxDuration.
			// Or if the satellite rises but doesn't set within the remaining maxDuration.
			// We can choose to stop, or log a warning and continue searching from riseTime + coarseStep.
			// For simplicity, we stop here.
			break
		}

		// 3. Add the valid range to the list
		timeRange := TimeRange{
			From: riseTime,
			To:   setTime,
		}
		diff := timeRange.To.Sub(timeRange.From)
		// Only add if duration is meaningful (longer than precision)
		if diff > precision {
			timeRange.Difference = diff.String()
			timeRangeList = append(timeRangeList, timeRange)
		} else {
			// If rise/set are too close, it might be a glitch or extremely short pass.
			// Skip it and continue searching.
		}

		// 4. Update currentTime to search for the *next* pass after the current one ends
		currentTime = setTime.Add(precision) // Start searching just after the set time
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
