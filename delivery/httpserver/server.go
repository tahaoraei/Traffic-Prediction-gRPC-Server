package httpserver

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"timeMachine/service/timeservice"
)

type Config struct {
	Port int
}

type Server struct {
	config  Config
	timeSvc timeservice.Service
	Router  *echo.Echo
}

func New(config Config, timeSvc timeservice.Service) Server {
	return Server{
		config:  config,
		timeSvc: timeSvc,
		Router:  echo.New(),
	}
}

func (s Server) Serve() {
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.Logger())

	s.Router.GET("/health/live", s.health)
	s.Router.GET("/health/ready", s.health)

	// Start server
	address := fmt.Sprintf(":%d", s.config.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}

func (s Server) health(c echo.Context) error {
	return c.JSON(http.StatusOK, "live and ready")
}
