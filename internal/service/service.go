package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BabyLev/Umka-1/internal/clients/r4uab"
	"github.com/BabyLev/Umka-1/internal/storage"
	"github.com/BabyLev/Umka-1/internal/types"
	"github.com/BabyLev/Umka-1/satellite"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	Storage     *storage.Storage
	r4uabClient *r4uab.Client
}

func New(storage *storage.Storage, r *r4uab.Client) *Service {
	return &Service{
		Storage:     storage,
		r4uabClient: r,
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

	sat, err := s.Storage.GetSatellite(int(req.SatelliteID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.Storage.GetSatellite: %w", err).Error()))
		return
	}

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

	sat, err := s.Storage.GetSatellite(int(req.SatelliteID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.Storage.GetSatellite: %w", err).Error()))
		return
	}

	var t time.Time

	if req.Timestamp == nil {
		t = time.Now().UTC()
	} else {
		t = time.Unix(*req.Timestamp, 0)
	}

	obsLoc, err := s.Storage.GetLocation(int(req.ObserverPositionID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.Storage.GetLocation: %w", err).Error()))
		return
	}

	coords := satellite.ObserverCoords{
		Lon: obsLoc.Location.Lon,
		Lat: obsLoc.Location.Lat,
		Alt: obsLoc.Location.Alt,
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

	sat, err := s.Storage.GetSatellite(int(req.SatelliteID))
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Errorf("s.Storage.GetSatellite: %w", err).Error()))
		return
	}

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
		err := s.Storage.DeleteSatellite(i)
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

	sats := s.Storage.FindSatellite(req.Name)

	var res FindSatelliteResponse

	res.Satellites = make(map[int]Satellite, len(sats))

	for k, s := range sats {
		satellite := Satellite{
			Line1: s.GetLine1(),
			Line2: s.GetLine2(),
			Name:  s.Name,
		}
		res.Satellites[k] = satellite
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling: %w", err).Error()))
		return
	}

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

	sat, err := s.Storage.GetSatellite(idInt)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.Storage.DeleteSatellite: %w", err).Error()))
		return
	}

	res.Line1 = sat.GetLine1()
	res.Line2 = sat.GetLine2()
	res.Name = sat.Name

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

	var sat satellite.Satellite

	if req.Satellite.NoradID == nil {
		sat = satellite.New(req.Satellite.Line1, req.Satellite.Line2)
	} else {
		updatedSatInfo, err := s.r4uabClient.GetSatelliteInfo(r.Context(), *req.Satellite.NoradID)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("r4uabClient.GetSatelliteInfo: %w", err).Error()))
			return
		}

		sat = satellite.New(updatedSatInfo.Line1, updatedSatInfo.Line2)
	}

	// надо проверить, передал ли клиент в запросе norad id.
	// если передал, то перезаписываем в структуре нужного спутника в хранилище
	if req.Satellite.NoradID == nil {
		err = s.Storage.UpdateSatellite(req.SatelliteID, types.Satellite{
			Satellite: sat,
			Name:      req.Satellite.Name,
			NoradID:   nil,
		})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("s.Storage.UpdateSatellite: %w", err).Error()))
			return
		}
	} else {
		err = s.Storage.UpdateSatellite(req.SatelliteID, types.Satellite{
			Satellite: sat,
			Name:      req.Satellite.Name,
			NoradID:   req.Satellite.NoradID,
		})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("s.Storage.UpdateSatellite: %w", err).Error()))
			return
		}
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

	var sat satellite.Satellite

	if req.NoradID == nil {
		sat = satellite.New(req.Line1, req.Line2)
	} else {
		updatedSatInfo, err := s.r4uabClient.GetSatelliteInfo(r.Context(), *req.NoradID)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Errorf("r4uabClient.GetSatelliteInfo: %w", err).Error()))
			return
		}

		sat = satellite.New(updatedSatInfo.Line1, updatedSatInfo.Line2)
	}

	satID := s.Storage.AddSatellite(types.Satellite{
		Satellite: sat,
		Name:      req.Name,
		NoradID:   req.NoradID,
	})

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

	locID := s.Storage.AddLocation(req.ObserverLocation)

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
