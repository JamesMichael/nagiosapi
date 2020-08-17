package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Server provides the main HTTP API.
type Server struct {
	addr            string
	log             *zap.Logger
	statusHandler   *StatusHandler
	commandsHandler *CommandsHandler
}

type ServerOpt func(s *Server) error

func NewServer(opts ...ServerOpt) (*Server, error) {
	s := &Server{
		addr: viper.GetString("api.addr"),
		log:  zap.NewNop(),
	}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// WithAddr can be used to customise the listen address.
func WithAddr(addr string) ServerOpt {
	return func(s *Server) error {
		s.addr = addr
		return nil
	}
}

func WithCommandsHandler(h *CommandsHandler) ServerOpt {
	return func(s *Server) error {
		s.commandsHandler = h
		return nil
	}
}

// WithLog can be used to pass in a Logger.
func WithLog(l *zap.Logger) ServerOpt {
	return func(s *Server) error {
		s.log = l
		return nil
	}
}

// WithStatusHandler can be used to set the handler for /api/status.
func WithStatusHandler(sh *StatusHandler) ServerOpt {
	return func(s *Server) error {
		s.statusHandler = sh
		return nil
	}
}

func (s *Server) ServeHTTP() {
	router := chi.NewRouter()

	cors := buildCORSMiddleware()
	if cors != nil {
		router.Use(cors.Handler)
	}

	auth := buildBasicAuthMiddleware()
	if auth != nil {
		router.Use(auth)
	}

	router.Use(
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/v1", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Route("/status", func(r chi.Router) {
				r.Get("/{host}/{service}", s.statusHandler.GetServiceStatus)
				r.Post("/", s.statusHandler.GetServiceStatusMulti)
			})
			r.Route("/submit", func(r chi.Router) {
				r.Post("/", s.commandsHandler.SubmitResult)
			})
		})
	})

	s.log.Info("starting HTTP API",
		zap.String("addr", s.addr),
	)
	if err := http.ListenAndServe(s.addr, router); err != nil {
		s.log.Fatal("unexpected server failure",
			zap.Error(err))
	}
}

func buildBasicAuthMiddleware() func(next http.Handler) http.Handler {
	if !viper.GetBool("basic_auth.enabled") {
		return nil
	}

	viper.SetDefault("basic_auth.realm", "nagios-api")
	viper.SetDefault("basic_auth.users", nil)

	return middleware.BasicAuth(viper.GetString("basic_auth.realm"), viper.GetStringMapString("basic_auth.users"))
}

func buildCORSMiddleware() *cors.Cors {
	if !viper.GetBool("cors.enabled") {
		return nil
	}

	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"})
	viper.SetDefault("cors.allow_credentials", false)
	viper.SetDefault("cors.max_age", 300)

	return cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice("cors.allowed_origins"),
		AllowedMethods:   viper.GetStringSlice("cors.allowed_methods"),
		AllowedHeaders:   viper.GetStringSlice("cors.allowed_headers"),
		AllowCredentials: viper.GetBool("cors.allow_credentials"),
		MaxAge:           viper.GetInt("cors.max_age"),
	})
}
