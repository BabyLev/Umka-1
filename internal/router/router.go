package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BabyLev/Umka-1/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Relative path to the frontend build output directory from the project root
const staticDirRelative string = "web/dist"

func SetupRouter(service *service.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// --- API routes --- Define API routes first
	router.Route("/calculate", func(r chi.Router) {
		r.Post("/", service.Calculate)
	})
	router.Route("/look-angles", func(r chi.Router) {
		r.Post("/", service.LookAngles)
	})
	router.Route("/time-ranges", func(r chi.Router) {
		r.Post("/", service.VisibleTimeRange)
	})
	router.Route("/satellite", func(r chi.Router) {
		r.Put("/", service.AddSatellite)
		r.Post("/", service.FindSatellite) // Keep POST for find as per service/readme
		r.Patch("/", service.UpdateSatellite)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", service.GetSatellite)
			r.Delete("/", service.DeleteSatellite)
		})
	})
	router.Route("/location", func(r chi.Router) {
		r.Put("/", service.AddLocation)
		r.Post("/", service.FindLocation)
		r.Patch("/", service.UpdateLocation)
		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", service.DeleteLocation)
			r.Get("/", service.GetLocation)
		})
	})

	// --- Static file serving for Vue SPA ---
	workDir, _ := os.Getwd()
	staticPath := filepath.Join(workDir, staticDirRelative)
	staticRoot := http.Dir(staticPath)

	// This function handles serving static files and the SPA index.html fallback.
	FileServer(router, "/", staticRoot)

	return router
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// It includes SPA fallback logic: if a requested file isn't found,
// it serves the index.html file.
func FileServer(r chi.Router, publicPath string, root http.FileSystem) {
	if strings.ContainsAny(publicPath, ":*") {
		panic("FileServer does not permit URL parameters.")
	}

	// Define the handler for serving files
	fs := http.StripPrefix(publicPath, http.FileServer(root))

	// Ensure the public path ends with a slash for prefix matching
	if publicPath != "/" && publicPath[len(publicPath)-1] != '/' {
		r.Get(publicPath, http.RedirectHandler(publicPath+"/", http.StatusMovedPermanently).ServeHTTP)
		publicPath += "/"
	}

	// Define the route for all paths under the public path
	r.Get(publicPath+"*", func(w http.ResponseWriter, r *http.Request) {
		// Get the requested file path relative to the static root
		requestedFilePath := strings.TrimPrefix(r.URL.Path, publicPath)
		// Try to open the file
		f, err := root.Open(requestedFilePath)
		if os.IsNotExist(err) {
			// File not found, serve index.html for SPA routing
			indexPath := filepath.Join(strings.TrimPrefix(publicPath, "/"), "index.html")
			http.ServeFile(w, r, filepath.Join(string(root.(http.Dir)), indexPath))
			return
		}
		if err != nil {
			// Other error (e.g., permission denied)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Serve the found file
		fs.ServeHTTP(w, r)
	})
}
