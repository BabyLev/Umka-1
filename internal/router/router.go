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

	router.Post("/calculate/", service.Calculate)             // возвращает длину, широту, долготу, адрес спутника на карте
	router.Post("/look-angles/", service.LookAngles)          // возвращает азимут, элевацию, диапазон спутника
	router.Post("/time-ranges/", service.VisibleTimeRange)    // высчитывает временные диапазоны, когда видно спутник над заданной точкой в нужном количестве
	router.Put("/satellite/", service.AddSatellite)           // добавляет переданный спутник в хранилище и возращает его ID
	router.Delete("/satellite/{id}", service.DeleteSatellite) // удаляет спутник из хранилища по ID
	router.Post("/satellite/", service.FindSatellite)         // возвращает все спутники по переданному имени
	router.Get("/satellite/{id}", service.GetSatellite)       // возвращает спутник по переданному ID
	router.Patch("/satellite/", service.UpdateSatellite)      // принимает ID и новые данные спутника. Изменяет старые переменные на новые
	router.Put("/satellite/", service.AddLocation)            // добавляет переданную локацию наблюдателя в хранилище и возращает ID
	// "/satellite/{id}"
	return router
}
