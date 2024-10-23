package storage

import (
	"fmt"
	"strings"

	"github.com/BabyLev/Umka-1/satellite"
)

// пакет нужен для хранения информации о спутниках (и не только)
// информацию храним в памяти (то есть, до окончания работы программы)

type Storage struct {
	Satellites []satellite.Satellite
	// SatellitesIndex map[string]int // key - name of sat, value - index in Satellites slice
	Locations []satellite.LookAngles
}

func New() *Storage {
	return &Storage{}
}

// функция вернет индекс в слайсе, под которым сохранен новый добавленный спутник
func (s *Storage) AddSatellite(sat satellite.Satellite) int {
	s.Satellites = append(s.Satellites, sat)

	return len(s.Satellites) - 1
}

func (s *Storage) GetSatellite(num int) (satellite.Satellite, error) {
	if num < 0 || num > len(s.Satellites)-1 {
		return satellite.Satellite{}, fmt.Errorf("нет спутника под таким ID")
	}

	return s.Satellites[num], nil
}

func (s *Storage) DeleteSatellite(num int) error {
	if num < 0 || num > len(s.Satellites)-1 {
		return fmt.Errorf("нет спутника под таким ID")
	}

	s.Satellites = append(s.Satellites[:num], s.Satellites[num+1:]...)

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
