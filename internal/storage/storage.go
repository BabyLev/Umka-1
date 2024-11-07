package storage

import (
	"fmt"
	"strings"

	"github.com/BabyLev/Umka-1/satellite"
)

// пакет нужен для хранения информации о спутниках (и не только)
// информацию храним в памяти (то есть, до окончания работы программы)

type Storage struct {
	satellites map[int]satellite.Satellite
	lastSatID  int
	locations  map[int]satellite.ObserverCoords
	lastLocID  int
}

// type SatellitesID struct {
// 	ID int
// }

func New() *Storage {
	sat := make(map[int]satellite.Satellite)
	loc := make(map[int]satellite.ObserverCoords)
	storage := Storage{
		satellites: sat,
		locations:  loc,
	}

	return &storage
}

// функция вернет индекс в слайсе, под которым сохранен новый добавленный спутник
func (s *Storage) AddSatellite(sat satellite.Satellite) int {
	s.lastSatID++
	s.satellites[s.lastSatID] = sat

	return s.lastSatID
}

func (s *Storage) GetSatellite(id int) (satellite.Satellite, error) {
	if _, ok := s.satellites[id]; !ok {
		return satellite.Satellite{}, fmt.Errorf("нет спутника под таким ID")
	}

	return s.satellites[id], nil
}

func (s *Storage) DeleteSatellite(id int) error {
	if _, ok := s.satellites[id]; !ok {
		return fmt.Errorf("нет спутника под таким ID")
	}

	delete(s.satellites, id)

	return nil
}

func (s *Storage) UpdateSatellite(id int, sat satellite.Satellite) error {
	if _, ok := s.satellites[id]; !ok {
		return fmt.Errorf("нет спутника под таким ID")
	}

	s.satellites[id] = sat

	return nil
}

// FindSatellite(name string) (satellite., error)
// for ex: sat name = umka, input name - um
func (s *Storage) FindSatellite(name string) map[int]satellite.Satellite {
	satellites := make(map[int]satellite.Satellite, 0)

	for i, sat := range s.satellites {
		if strings.Contains(sat.GetName(), name) {
			satellites[i] = sat
		}
	}

	return satellites
}

func (s *Storage) GetLocation(ID int) (satellite.ObserverCoords, error) {
	loc, ok := s.locations[ID]
	if !ok {
		return satellite.ObserverCoords{}, fmt.Errorf("нет локации под таким ID")
	}

	return loc, nil
}

func (s *Storage) AddLocation(loc satellite.ObserverCoords) int {
	s.lastLocID++
	s.locations[s.lastLocID] = loc

	return s.lastLocID
}
