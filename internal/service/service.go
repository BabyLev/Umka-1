package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BabyLev/Umka-1/satellite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// routes // маршруты
	// маршрут для главной страницы
	router.Get("/", MainPage) // "/" - path, root path // корневой путь

	router.Get("/calculate/", Calculate)
	router.Post("/look-angles/", LookAngles)

	return router
}

// ручка для обработки главной страницы
func MainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}

func Calculate(w http.ResponseWriter, r *http.Request) {
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
func LookAngles(w http.ResponseWriter, r *http.Request) {
	var req LookAnglesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Errorf("ошибка декодирования запроса: %w", err).Error()))
		return
	}

	sat := satellite.New(req.Line1, req.Line2, req.SatName)

	t := time.Unix(req.Timestamp, 0)

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