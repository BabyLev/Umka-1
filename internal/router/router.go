package router

import (
	"net/http"
	"path"

	"github.com/BabyLev/Umka-1/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const staticDir string = "../static/"

func SetupRouter(service *service.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// routes // маршруты
	router.Get("/", ServeStatic)

	// расчет параметров спутника
	router.Post("/calculate/", service.Calculate)          // возвращает длину, широту, долготу, адрес спутника на карте
	router.Post("/look-angles/", service.LookAngles)       // возвращает азимут, элевацию, диапазон спутника
	router.Post("/time-ranges/", service.VisibleTimeRange) // высчитывает временные диапазоны, когда видно спутник над заданной точкой в нужном количестве

	// управление сохраненными спутниками
	router.Put("/satellite/", service.AddSatellite)           // добавляет переданный спутник в хранилище и возращает его ID
	router.Delete("/satellite/{id}", service.DeleteSatellite) // удаляет спутник из хранилища по ID
	router.Post("/satellite/", service.FindSatellite)         // возвращает все спутники по переданному имени
	router.Get("/satellite/{id}", service.GetSatellite)       // возвращает спутник по переданному ID
	router.Patch("/satellite/", service.UpdateSatellite)      // принимает ID и новые данные спутника. Изменяет старые переменные на новые

	// управление сохраненными локациями
	router.Put("/location/", service.AddLocation)           // добавляет переданную локацию наблюдателя в хранилище и возращает ID
	router.Delete("/location/{id}", service.DeleteLocation) // удаляет локацию наблюдателя из хранилища по ID
	router.Post("/location/", service.FindLocation)         // возвращает все локации наблюдателя по переданному имени
	router.Patch("/location/", service.UpdateLocation)      // принимает ID и новые данные локации наблюдателя. Изменяет старые переменные на новые
	// "/satellite/{id}"
	return router
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(staticDir, "index.html"))
}
