package storage

import (
	"fmt"
	"strings"

	"github.com/BabyLev/Umka-1/satellite"
)

// пакет нужен для хранения информации о спутниках (и не только)
// информацию храним в памяти (то есть, до окончания работы программы)

type Storage struct {
	Satellites map[int]satellite.Satellite
	LastSatID  int
	// SatellitesIndex map[string]int // key - name of sat, value - index in Satellites slice
	Locations []satellite.LookAngles
}

// type SatellitesID struct {
// 	ID int
// }

func New() *Storage {
	sat := make(map[int]satellite.Satellite)

	return &Storage{
		Satellites: sat,
	}
}

// функция вернет индекс в слайсе, под которым сохранен новый добавленный спутник
func (s *Storage) AddSatellite(sat satellite.Satellite) int {
	s.LastSatID++
	s.Satellites[s.LastSatID] = sat

	return s.LastSatID
}

func (s *Storage) GetSatellite(id int) (satellite.Satellite, error) {
	if _, ok := s.Satellites[id]; !ok {
		return satellite.Satellite{}, fmt.Errorf("нет спутника под таким ID")
	}

	return s.Satellites[id], nil
}

func (s *Storage) DeleteSatellite(id int) error {
	if _, ok := s.Satellites[id]; !ok {
		return fmt.Errorf("нет спутника под таким ID")
	}

	delete(s.Satellites, id)

	return nil
}

func (s *Storage) UpdateSatellite(id int, sat satellite.Satellite) error {
	if _, ok := s.Satellites[id]; !ok {
		return fmt.Errorf("нет спутника под таким ID")
	}

	s.Satellites[id] = sat

	return nil
}

// FindSatellite(name string) (satellite., error)
// for ex: sat name = umka, input name - um
func (s *Storage) FindSatellite(name string) map[int]satellite.Satellite {
	satellites := make(map[int]satellite.Satellite, 0)

	for i, sat := range s.Satellites {
		if strings.Contains(sat.GetName(), name) {
			satellites[i] = sat
		}
	}

	return satellites
}
