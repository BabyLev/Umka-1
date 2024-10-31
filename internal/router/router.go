package router

import (
	"github.com/BabyLev/Umka-1/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(service *service.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// routes // маршруты
	// маршрут для главной страницы
	router.Get("/", service.MainPage) // "/" - path, root path // корневой путь

	router.Get("/calculate/", service.Calculate)
	router.Post("/look-angles/", service.LookAngles)
	router.Post("/time-ranges/", service.VisibleTimeRange)
	router.Put("/satellite/", service.AddSatellite)
	router.Delete("/satellite/{id}", service.DeleteSatellite)
	router.Post("/satellite/", service.FindSatellite)
	router.Get("/satellite/{id}", service.GetSatellite)
	router.Patch("/satellite/", service.UpdateSatellite)
	// "/satellite/{id}"
	return router
}
