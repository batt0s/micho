package api

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/batt0s/micho/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type API struct {
	Addr   string
	Router *chi.Mux
	Server http.Server
}

func (api *API) Init() error {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	accessLogFile, err := os.OpenFile(filepath.Join("logs", "access.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	mWriter := io.MultiWriter(os.Stdout, accessLogFile)
	logger := middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(mWriter, "", log.LstdFlags),
		NoColor: true,
	})
	r.Use(logger)

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 240))

	if err := logging.Init(); err != nil {
		return err
	}

	r.Route("/api", func(ar chi.Router) {
		ar.Use(APIKeyMiddleware)
		ar.Post("/deploy", DeployPyMenuHandler)
		ar.Post("/upgrade", UpgradePyMenuHandler)
		ar.Delete("/uninstall/{slug}", UninstallPyMenuHandler)
		ar.Get("/status/{slug}", StatusPyMenuHandler)
	})

	var host, port string
	host = os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port = os.Getenv("PORT")
	if port == "" {
		port = "3578"
	}

	api.Router = r
	api.Addr = host + ":" + port
	api.Server = http.Server{
		Addr:    api.Addr,
		Handler: api.Router,
	}
	return nil
}
