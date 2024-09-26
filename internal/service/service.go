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

	w.Write([]byte(res))
}
