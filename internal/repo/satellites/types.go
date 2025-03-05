package satellites

type Satellite struct {
	ID      int
	SatName string
	NoradID *int
	Line1   string
	Line2   string
}

type FilterSatellite struct {
	IDs      []int
	SatName  *string
	NoradIDs []int
}
