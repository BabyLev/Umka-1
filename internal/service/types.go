package service

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
	Lon               float64 `json:"lat"`
	Lat               float64 `json:"lon"`
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
	Name string `json:"name"`
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
	NoradID *int   `json:"noradId"`
}
