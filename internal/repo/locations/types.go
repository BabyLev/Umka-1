package locations

type Location struct {
	ID    int
	Name  string
	Point Point
}

type Point struct {
	Lon float64 // долгота
	Lat float64 // широта
	Alt float64 // высота
}

type FilterLocation struct {
	Name *string
}
