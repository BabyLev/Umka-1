package service

type LookAnglesRequest struct {
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	SatName   string `json:"satName"`
	Timestamp *int64 `json:"timestamp"`
	// Координаты наблюдателя
	Lon float64 `json:"lat"`
	Lat float64 `json:"lon"`
	Alt float64 `json:"alt"` // км
}

type VisibleTimeRangeRequest struct {
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	SatName   string `json:"satName"`
	Timestamp *int64 `json:"timestamp"`
	// Координаты наблюдателя
	Lon               float64 `json:"lat"`
	Lat               float64 `json:"lon"`
	Alt               float64 `json:"alt"` // км
	CountOfTimeRanges *int    `json:"countOfTimeRanges"`
}

type AddSatelliteRequest struct {
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Name  string `json:"name"`
}

type AddSatelliteResponse struct {
	ID int64 `json:"id"`
}

type FindSatelliteRequest struct {
	Name string `json:"name"`
}

type FindSatelliteResponse struct {
	Satellites map[int]Satellite `json:"satellites"` // int - id cпутника
}

type UpdateSatelliteRequest struct {
	Satellite Satellite `json:"satellite"`
	ID        int       `json:"id"`
}

type Satellite struct {
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Name  string `json:"name"`
}
