package httpserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trinity/api/router"
	"trinity/configs"
	"trinity/pkg/app"
)

type Server struct {
	appConfig *app.Config
	cfg       *configs.Config
	httpSrv   *http.Server
	router    *gin.Engine
}

func NewHTTPServer(cfg *configs.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}
}

func NewServer(appConfig *app.Config, cfg *configs.Config) *Server {
	router := router.NewRouter(appConfig, cfg)
	httpSrv := NewHTTPServer(cfg, router)
	return &Server{
		cfg:     cfg,
		httpSrv: httpSrv,
		router:  router,
	}
}

func (s *Server) Start() error {
	logrus.Info("Starting server on ", s.httpSrv.Addr)
	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server failed to start: %v", err)
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	logrus.Info("Shutting down server...")
	return s.httpSrv.Shutdown(ctx)
}

func (s *Server) StartWithGracefulShutdown() {
	if err := s.Start(); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		logrus.Fatalf("Failed to stop server gracefully: %v", err)
	}
	logrus.Info("Server stopped")
}
