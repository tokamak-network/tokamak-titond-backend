package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/api"
)

type Config struct {
	Host string
	Port string
}

type HTTPServer struct {
	cfg  *Config
	apis *api.TitondAPI
}

func NewHTTPServer(cfg *Config, apis *api.TitondAPI) *HTTPServer {
	return &HTTPServer{
		cfg,
		apis,
	}
}

func (s *HTTPServer) Run() {
	r := gin.Default()

	networkRouter := r.Group("/api/networks")
	s.NewNetworkRouter(networkRouter)

	resourceRouter := r.Group("/api/resources")
	s.NewResourceRouter(resourceRouter)

	r.Run(s.cfg.Host + ":" + s.cfg.Port)
}
