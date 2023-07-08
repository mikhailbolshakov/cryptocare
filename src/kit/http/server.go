package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mikhailbolshakov/cryptocare/src/kit/goroutine"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"github.com/rs/cors"
	"net/http"
	"time"
)

type Cors struct {
	AllowedHeaders []string
	AllowedOrigins []string
	AllowedMethods []string
	Debug          bool
}

type Config struct {
	Port                 string
	Cors                 *Cors
	Trace                bool
	WriteTimeoutSec      int `config:"write-timeout-sec"`
	ReadTimeoutSec       int `config:"read-timeout-sec"`
	ReadBufferSizeBytes  int `config:"read-buffer-size-bytes"`
	WriteBufferSizeBytes int `config:"write-buffer-size-bytes"`
}

// Server represents HTTP server
type Server struct {
	Srv        *http.Server        // Srv - internal server
	RootRouter *mux.Router         // RootRouter - root router
	WsUpgrader *websocket.Upgrader // WsUpgrader - websocket upgrader
	logger     log.CLoggerFunc     // logger
}

type RouteSetter interface {
	Set() error
}

type WsUpgrader interface {
	Set(router *mux.Router, upgrader *websocket.Upgrader)
}

// getOptions getting cors options preconfigured
func getOptions(cfg *Config) cors.Options {
	if cfg.Cors == nil {
		return cors.Options{
			AllowCredentials: true,
		}
	}

	return cors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		AllowCredentials: true,
		Debug:            cfg.Cors.Debug,
	}
}

func NewHttpServer(cfg *Config, logger log.CLoggerFunc) *Server {
	r := mux.NewRouter()

	corsHandler := cors.New(getOptions(cfg)).Handler(r)

	s := &Server{
		Srv: &http.Server{
			Addr:         fmt.Sprintf(":%s", cfg.Port),
			Handler:      corsHandler,
			WriteTimeout: time.Duration(cfg.WriteTimeoutSec) * time.Second,
			ReadTimeout:  time.Duration(cfg.ReadTimeoutSec) * time.Second,
		},
		WsUpgrader: &websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufferSizeBytes,
			WriteBufferSize: cfg.WriteBufferSizeBytes,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		logger: logger,
	}
	if cfg.Trace {
		r.Use(s.loggingMiddleware)
	}
	s.RootRouter = r

	return s
}

func (s *Server) SetWsUpgrader(upgradeSetter WsUpgrader) {
	upgradeSetter.Set(s.RootRouter, s.WsUpgrader)
}

func (s *Server) Listen() {

	goroutine.New().
		WithLoggerFn(s.logger).
		WithRetry(goroutine.Unrestricted).
		Cmp("http-server").
		Mth("listen").
		Go(context.Background(),
			func() {
				l := s.logger().Pr("http").Cmp("server").Mth("listen").F(log.FF{"url": s.Srv.Addr})
				l.Inf("start listening")
			start:
				if err := s.Srv.ListenAndServe(); err != nil {
					if err != http.ErrServerClosed {
						l.E(ErrHttpSrvListen(err)).St().Err()
						time.Sleep(time.Second * 5)
						goto start
					} else {
						l.Dbg("server closed")
					}
					return
				}
			})
}

func (s *Server) Close() {
	_ = s.Srv.Close()
}
