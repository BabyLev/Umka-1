package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BabyLev/Umka-1/internal/clients/r4uab"
	locationsRepo "github.com/BabyLev/Umka-1/internal/repo/locations"
	satellitesRepo "github.com/BabyLev/Umka-1/internal/repo/satellites"
	"github.com/BabyLev/Umka-1/internal/types"
	"github.com/BabyLev/Umka-1/satellite"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	repoSats    *satellitesRepo.Repo
	repoLocs    *locationsRepo.Repo
	r4uabClient *r4uab.Client
}

func New(rClient *r4uab.Client, repoSats *satellitesRepo.Repo, repoLocs *locationsRepo.Repo) *Service {
	return &Service{
		r4uabClient: rClient,
		repoSats:    repoSats,
		repoLocs:    repoLocs,
	}
}

// ручка для обработки главной страницы
func (s *Service) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Добро пожаловать на сервис управления спутниками"))
}

func (s *Service) Calculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	var t time.Time

	if req.Timestamp != nil {
		t = time.Unix(int64(*req.Timestamp), 0)
	} else {
		t = time.Now()
	}

	satRepo, err := s.repoSats.GetSatellite(r.Context(), int(req.SatelliteID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.repo.GetSatellite: %w", err).Error()))
		return
	}

	sat := satellite.New(satRepo.Line1, satRepo.Line2)

	satCoords, err := sat.Calculate(t.UTC())
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("ошибка при подсчёте координат: %w", err).Error()))
		return
	}

	res, err := json.Marshal(satCoords)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling coords: %w", err).Error()))
		return
	}

	w.Write(res)
}

// POST /look_angles
// Example request
//
//	{
//	    "satelliteId": 3,
//	    "timestamp": 1727978254,
//	    "lat": 55.43025412211996,
//	    "lon": 37.51934842793972,
//	    "alt": 0.151
//	}
func (s *Service) LookAngles(w http.ResponseWriter, r *http.Request) {
	var req LookAnglesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	satRepo, err := s.repoSats.GetSatellite(r.Context(), int(req.SatelliteID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.repo.GetSatellite: %w", err).Error()))
		return
	}

	sat := satellite.New(satRepo.Line1, satRepo.Line2)

	var t time.Time

	if req.Timestamp == nil {
		t = time.Now().UTC()
	} else {
		t = time.Unix(*req.Timestamp, 0)
	}

	obsLoc, err := s.repoLocs.GetLocation(r.Context(), int(req.ObserverPositionID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.Storage.GetLocation: %w", err).Error()))
		return
	}

	coords := satellite.ObserverCoords{
		Lon: obsLoc.Point.Lon,
		Lat: obsLoc.Point.Lat,
		Alt: obsLoc.Point.Alt,
	}

	lookAngles := sat.LookAngles(t, coords)

	res, err := json.Marshal(lookAngles)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling coords: %w", err).Error()))
		return
	}

	w.Write(res)
}

func (s *Service) VisibleTimeRange(w http.ResponseWriter, r *http.Request) {
	var req VisibleTimeRangeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	satRepo, err := s.repoSats.GetSatellite(r.Context(), int(req.SatelliteID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.repo.GetSatellite: %w", err).Error()))
		return
	}

	sat := satellite.New(satRepo.Line1, satRepo.Line2)

	var t time.Time

	if req.Timestamp == nil {
		t = time.Now().UTC()
	} else {
		t = time.Unix(*req.Timestamp, 0)
	}

	countOfTimeRanges := 1
	if req.CountOfTimeRanges != nil {
		countOfTimeRanges = *req.CountOfTimeRanges
	}

	timeRanges := sat.VisibleTimeRange(t, satellite.ObserverCoords{
		Lon: req.Lon,
		Lat: req.Lat,
		Alt: req.Alt,
	}, countOfTimeRanges)

	res, err := json.Marshal(timeRanges)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling time ranges: %w", err).Error()))
		return
	}

	w.Write(res)
}

// HTTP Method: DELETE
// URL Param: satellite id
// http://localhost/satellite/{id}
// Example: DELETE http://localhost/satellite/123
func (s *Service) DeleteSatellite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // "123"

	if i, err := strconv.Atoi(id); err == nil {
		err = s.repoSats.DeleteSatellite(r.Context(), i)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("s.Storage.DeleteSatellite: %w", err).Error()))
			return
		}
	} else {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("не удалось преобразовать ID к целому числу: %w", err).Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("спутник успешно удалился id = %s", id)))
}

// HTTP Method: POST
// Body: {"name": "something"}
// http://localhost/satellite
func (s *Service) FindSatellite(w http.ResponseWriter, r *http.Request) {
	var req FindSatelliteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	filter := satellitesRepo.FilterSatellite{
		IDs:      req.IDs,
		NoradIDs: req.NoradIDs,
		SatName:  req.Name,
	}

	sats, err := s.repoSats.FindSatellite(r.Context(), filter)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.repo.FindSatellite: %w", err).Error()))
		return
	}

	var res FindSatelliteResponse

	res.Satellites = make(map[int]Satellite, len(sats))

	for id, s := range sats {
		satellite := Satellite{
			Line1:   s.Line1,
			Line2:   s.Line2,
			Name:    s.SatName,
			NoradID: s.NoradID,
		}
		res.Satellites[id] = satellite
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling: %w", err).Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resJSON)
}

func (s *Service) GetSatellite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var res Satellite

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("ID невозможно преобразовать в число: %w", err).Error()))
	}

	satRepo, err := s.repoSats.GetSatellite(r.Context(), idInt)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.Repo.GetSatellite: %w", err).Error()))
		return
	}

	res.Line1 = satRepo.Line1
	res.Line2 = satRepo.Line2
	res.Name = satRepo.SatName
	res.NoradID = satRepo.NoradID

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling: %w", err).Error()))
		return
	}

	w.Write(resJSON)
}

func (s *Service) UpdateSatellite(w http.ResponseWriter, r *http.Request) {
	var req UpdateSatelliteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	if req.Satellite.NoradID != nil {
		updatedSatInfo, err := s.r4uabClient.GetSatelliteInfo(r.Context(), *req.Satellite.NoradID)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("r4uabClient.GetSatelliteInfo: %w", err).Error()))
			return
		}

		req.Satellite.Line1 = updatedSatInfo.Line1
		req.Satellite.Line2 = updatedSatInfo.Line2
	}

	satRepo := satellitesRepo.Satellite{
		ID:      req.SatelliteID,
		SatName: req.Satellite.Name,
		NoradID: req.Satellite.NoradID,
		Line1:   req.Satellite.Line1,
		Line2:   req.Satellite.Line2,
	}
	err = s.repoSats.UpdateSatellite(r.Context(), satRepo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.repo.UpdateSatellite: %w", err).Error()))
		return
	}

	w.WriteHeader(200)
}

func (s *Service) AddSatellite(w http.ResponseWriter, r *http.Request) {
	var req AddSatelliteRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	if req.NoradID != nil {
		updatedSatInfo, err := s.r4uabClient.GetSatelliteInfo(r.Context(), *req.NoradID)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("r4uabClient.GetSatelliteInfo: %w", err).Error()))
			return
		}

		req.Line1 = updatedSatInfo.Line1
		req.Line2 = updatedSatInfo.Line2
	}

	satID, err := s.repoSats.CreateSatellite(r.Context(), satellitesRepo.Satellite{
		SatName: req.Name,
		NoradID: req.NoradID,
		Line1:   req.Line1,
		Line2:   req.Line2,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("repo.CreateSatellite: %w", err).Error()))
		return
	}

	res := AddSatelliteResponse{
		SatelliteID: int64(satID),
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling: %w", err).Error()))
		return
	}

	w.Write(resJSON)
}

func (s *Service) AddLocation(w http.ResponseWriter, r *http.Request) {
	var req AddLocationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	obs := locationsRepo.Location{
		Name: req.ObserverLocation.Name,
		Point: locationsRepo.Point{
			Lon: req.ObserverLocation.Location.Lon,
			Lat: req.Location.Lat,
			Alt: req.Location.Alt,
		},
	}

	locID, err := s.repoLocs.CreateLocation(r.Context(), obs)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.repoLocs.CreateLocation: %w", err).Error()))
		return
	}

	res := AddLocationResponse{
		ID: locID,
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling: %w", err).Error()))
		return
	}

	w.Write(resJSON)
}

func (s *Service) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // "123"

	if i, err := strconv.Atoi(id); err == nil {
		err := s.repoLocs.DeleteLocation(r.Context(), i)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("s.repoLocs.DeleteLocation: %w", err).Error()))
			return
		}
	} else {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("не удалось преобразовать ID к целому числу: %w", err).Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("локация успешно удалилась id = %s", id)))
}

func (s *Service) FindLocation(w http.ResponseWriter, r *http.Request) {
	var req FindLocationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	filter := locationsRepo.FilterLocation{
		Name: &req.Name,
	}

	locs, err := s.repoLocs.FindLocation(r.Context(), filter)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.repoLocs.FindLocation: %w", err).Error()))
		return
	}

	resLocations := make(map[int]types.ObserverLocation) // key - location id in db

	for _, loc := range locs {
		resLocations[loc.ID] = types.ObserverLocation{
			Name: loc.Name,
			Location: types.Location{
				Lon: loc.Point.Lon,
				Lat: loc.Point.Lat,
				Alt: loc.Point.Alt,
			},
		}
	}

	res := FindLocationResponse{
		Locations: resLocations,
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling: %w", err).Error()))
		return
	}

	w.Write(resJSON)
}

func (s *Service) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	var req UpdateLocationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	resLocation := locationsRepo.Location{
		ID:   req.LocationID,
		Name: req.Location.Name,
		Point: locationsRepo.Point{
			Lon: req.Location.Location.Lon,
			Lat: req.Location.Location.Lat,
			Alt: req.Location.Location.Alt,
		},
	}

	err = s.repoLocs.UpdateLocation(r.Context(), resLocation)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.repoLocs.UpdateLocation: %w", err).Error()))
		return
	}

	w.WriteHeader(200)
}
