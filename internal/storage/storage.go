package storage

import (
	"fmt"
	"strings"

	"github.com/BabyLev/Umka-1/internal/types"
)

// пакет нужен для хранения информации о спутниках (и не только)
// информацию храним в памяти (то есть, до окончания работы программы)
// CRUD - Create/Read/Update/Delete

type Storage struct {
	satellites map[int]types.Satellite
	lastSatID  int
	locations  map[int]types.ObserverLocation
	lastLocID  int
}

// type SatellitesID struct {
// 	ID int
// }

func New() *Storage {
	sat := make(map[int]types.Satellite)
	loc := make(map[int]types.ObserverLocation)
	storage := Storage{
		satellites: sat,
		locations:  loc,
	}

	return &storage
}

// функция вернет индекс в слайсе, под которым сохранен новый добавленный спутник
func (s *Storage) AddSatellite(sat types.Satellite) int {
	s.lastSatID++
	s.satellites[s.lastSatID] = sat

	return s.lastSatID
}

func (s *Storage) GetSatellite(id int) (types.Satellite, error) {
	if _, ok := s.satellites[id]; !ok {
		return types.Satellite{}, fmt.Errorf("нет спутника под таким ID")
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

func (s *Storage) UpdateSatellite(id int, sat types.Satellite) error {
	if _, ok := s.satellites[id]; !ok {
		return fmt.Errorf("нет спутника под таким ID")
	}

	s.satellites[id] = sat

	return nil
}

// FindSatellite(name string) (satellite., error)
// for ex: sat name = umka, input name - um
func (s *Storage) FindSatellite(name string) map[int]types.Satellite {
	satellites := make(map[int]types.Satellite, 0)

	for i, sat := range s.satellites {
		if strings.Contains(sat.Name, name) {
			satellites[i] = sat
		}
	}

	return satellites
}

func (s *Storage) GetLocation(ID int) (types.ObserverLocation, error) {
	loc, ok := s.locations[ID]
	if !ok {
		return types.ObserverLocation{}, fmt.Errorf("нет локации под таким ID")
	}

	return loc, nil
}

func (s *Storage) AddLocation(loc types.ObserverLocation) int {
	s.lastLocID++
	s.locations[s.lastLocID] = loc

	return s.lastLocID
}

// TODO: Delete, Find, Update

func (s *Storage) DeleteLocation(ID int) error {
	_, ok := s.locations[ID]
	if !ok {
		return fmt.Errorf("нет локации под таким ID")
	}

	delete(s.locations, ID)

	return nil
}

func (s *Storage) FindLocation(name string) (map[int]types.ObserverLocation, error) {
	locals := make(map[int]types.ObserverLocation, 0)

	name = strings.ToLower(name)

	for k, v := range s.locations {
		if strings.Contains(strings.ToLower(v.Name), name) {
			locals[k] = v
		}
	}

	return locals, nil
}

func (s *Storage) UpdateLocation(ID int, loc types.ObserverLocation) error {
	_, ok := s.locations[ID]
	if !ok {
		return fmt.Errorf("нет локации под таким ID")
	}

	s.locations[ID] = loc

	return nil
}
