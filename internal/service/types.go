package service

import "github.com/BabyLev/Umka-1/internal/types"

type CalculateRequest struct {
	SatelliteID int64  `json:"satelliteId"` // id спутника из хранилища
	Timestamp   *int64 `json:"timestamp"`
}

type LookAnglesRequest struct {
	SatelliteID int64  `json:"satelliteId"` // id спутника из хранилища
	Timestamp   *int64 `json:"timestamp"`
	// Координаты наблюдателя
	ObserverPositionID int64 `json:"observerPositionId"`
}

type VisibleTimeRangeRequest struct {
	SatelliteID int64  `json:"satelliteId"` // id спутника из хранилища
	Timestamp   *int64 `json:"timestamp"`
	// Координаты наблюдателя
	Lon               float64 `json:"lon"`
	Lat               float64 `json:"lat"`
	Alt               float64 `json:"alt"` // км
	CountOfTimeRanges *int    `json:"countOfTimeRanges"`
}

type AddSatelliteRequest struct {
	Satellite
}

type AddSatelliteResponse struct {
	SatelliteID int64 `json:"satelliteId"`
}

type FindSatelliteRequest struct {
	IDs      []int   `json:"ids"`
	NoradIDs []int   `json:"noradIds"`
	Name     *string `json:"name"`
}

type FindSatelliteResponse struct {
	Satellites map[int]Satellite `json:"satellites"` // int - id cпутника в хранилище
}

type UpdateSatelliteRequest struct {
	Satellite   Satellite `json:"satellite"`
	SatelliteID int       `json:"satelliteId"` // id спутника в хранилище
}

type Satellite struct {
	Line1   string `json:"line1"`
	Line2   string `json:"line2"`
	Name    string `json:"name"`
	NoradID *int64 `json:"noradId"`
}

type AddLocationRequest struct {
	types.ObserverLocation
}

type AddLocationResponse struct {
	ID int `json:"observerLocationId"`
}

type FindLocationRequest struct {
	Name string `json:"name"`
}

type FindLocationResponse struct {
	Locations map[int]types.ObserverLocation `json:"locations"`
}

type UpdateLocationRequest struct {
	Location   types.ObserverLocation `json:"location"`
	LocationID int                    `json:"locationId"`
}
