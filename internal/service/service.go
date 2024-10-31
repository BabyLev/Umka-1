package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BabyLev/Umka-1/internal/storage"
	"github.com/BabyLev/Umka-1/satellite"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	Storage *storage.Storage
}

func New(storage *storage.Storage) *Service {
	return &Service{
		Storage: storage,
	}
}

// ручка для обработки главной страницы
func (s *Service) MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Добро пожаловать на сервис управления спутниками"))
}

func (s *Service) Calculate(w http.ResponseWriter, r *http.Request) {
	urlValues := r.URL.Query()

	line1 := urlValues.Get("line1")
	line2 := urlValues.Get("line2")
	name := urlValues.Get("name")
	timeStr := urlValues.Get("time")

	var t time.Time

	if timeStr != "" {
		ts, err := strconv.Atoi(timeStr)
		if err != nil {
			w.Write([]byte(fmt.Errorf("time: %w", err).Error()))
			return
		}

		t = time.Unix(int64(ts), 0)
	} else {
		t = time.Now()
	}

	sat := satellite.New(line1, line2, name)

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
//	    "line1": "1 57172U 23091G   24277.18824023  .00015270  00000-0  94531-3 0  9997",
//	    "line2": "2 57172  97.5992 328.2179 0017737 111.7890 248.5228 15.09757499 69775",
//	    "satName": "UMKA-(RS40S)",
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

	sat := satellite.New(req.Line1, req.Line2, req.SatName)

	var t time.Time

	if req.Timestamp == nil {
		t = time.Now().UTC()
	} else {
		t = time.Unix(*req.Timestamp, 0)
	}

	lookAngles := sat.LookAngles(t, satellite.ObserverCoords{
		Lon: req.Lon,
		Lat: req.Lat,
		Alt: req.Alt,
	})

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

	sat := satellite.New(req.Line1, req.Line2, req.SatName)

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

	// mapping: преобразование одного типа данных в другой
	// задача сделать из map[int]satellite.Satellite -> map[int]Satellite

	res.Satellites = make(map[int]Satellite, len(sats))

	for k, s := range sats {
		satellite := Satellite{
			Line1: s.GetLine1(),
			Line2: s.GetLine2(),
			Name:  s.GetName(),
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
	res.Name = sat.GetName()

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

	sat := satellite.New(req.Satellite.Line1, req.Satellite.Line2, req.Satellite.Name)
	err = s.Storage.UpdateSatellite(req.ID, sat)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("s.Storage.UpdateSatellite: %w", err).Error()))
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

	sat := satellite.New(req.Line1, req.Line2, req.Name)

	satID := s.Storage.AddSatellite(sat)

	res := AddSatelliteResponse{
		ID: int64(satID),
	}

	resJSON, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("error marshalling: %w", err).Error()))
		return
	}

	w.Write(resJSON)
}
