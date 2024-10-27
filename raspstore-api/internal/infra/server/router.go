package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/murilo-bracero/raspstore/file-service/internal/auth"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
)

const serviceBaseRoute = "/file-service"
const fileBaseRoute = serviceBaseRoute + "/v1/files"
const uploadRoute = serviceBaseRoute + "/v1/uploads"
const downloadRoute = serviceBaseRoute + "/v1/downloads/{fileId}"
const loginRoute = serviceBaseRoute + "/v1/login"
const registerRoute = serviceBaseRoute + "/v1/register"

type FilesRouter interface {
	MountRoutes() *chi.Mux
}

type filesRouter struct {
	config       *config.Config
	handler      *handler.Handler
	jwtValidator *validator.JWTValidator
}

func NewFilesRouter(
	config *config.Config,
	handler *handler.Handler,
	jwtValidator *validator.JWTValidator) FilesRouter {
	return &filesRouter{config: config, handler: handler, jwtValidator: jwtValidator}
}

func (fr *filesRouter) MountRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)

	// private routes
	router.Route("/", func(r chi.Router) {
		r.Use(auth.TokenMiddleware(fr.jwtValidator))

		r.Route(fileBaseRoute, func(r1 chi.Router) {
			r1.Get("/", fr.handler.ListFiles)
			r1.Get("/{id}", fr.handler.FindById)
			r1.Put("/{id}", fr.handler.Update)
			r1.Delete("/{id}", fr.handler.Delete)
		})

		r.Post(uploadRoute, fr.handler.Upload)
		r.Get(downloadRoute, fr.handler.Download)
	})

	// public routes
	router.Post(loginRoute, fr.handler.Authenticate)
	router.Post(registerRoute, fr.handler.Register)

	return router
}

func cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: add config field to add origins based on env variables
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
