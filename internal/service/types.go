package service

type LookAnglesRequest struct {
	Line1     string `json:"line1"`
	Line2     string `json:"line2"`
	SatName   string `json:"satName"`
	Timestamp int64  `json:"timestamp"`
	// Координаты наблюдателя
	Lon float64 `json:"lat"`
	Lat float64 `json:"lon"`
	Alt float64 `json:"alt"` // км
}
