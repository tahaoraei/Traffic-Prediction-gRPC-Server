package httpserver

import (
	"fmt"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"timeMachine/pkg/logger"
	"timeMachine/pkg/metric"
)

var log = logger.Get()

type Config struct {
	Port int
}

type Server struct {
	config Config
	Router *echo.Echo
}

func New(config Config) Server {
	return Server{
		config: config,
		Router: echo.New(),
	}
}

func (s Server) Serve() {
	if err := prometheus.Register(metric.ResponseHistogram); err != nil {
		log.Fatal().Msgf("can't register prometheus metric: ", err)
	}

	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.Logger())

	s.Router.GET("/health/live", s.health)
	s.Router.GET("/health/ready", s.health)
	s.Router.GET("/actuator/prometheus", echoprometheus.NewHandler())

	// Start server
	address := fmt.Sprintf(":%d", s.config.Port)
	log.Info().Str("server", "start server").Msgf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		log.Fatal().Msgf("router start error", err)
	}
}

func (s Server) health(c echo.Context) error {
	return c.JSON(http.StatusOK, "live and ready")
}
