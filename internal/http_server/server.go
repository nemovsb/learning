package http_server

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type ServerConfig struct {
	Port     string
	RTimeout uint
	WTimeout uint
}

type Server struct {
	server http.Server
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
	//return http.ListenAndServe(s.server.Addr, s.server.Handler)
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}

func NewServer(conf ServerConfig, handler http.Handler) *Server {
	fmt.Printf("Port : %s", fmt.Sprint(":"+conf.Port))
	serv := http.Server{
		Addr:    ":" + conf.Port,
		Handler: handler,
	}

	return &Server{server: serv}
}

type AppServer struct {
	Port        string
	RTimeout    uint
	WTimeout    uint
	RedisClient *redis.Client
}
