package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tokamak-network/tokamak-titond-backend/docs"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	networkRouter := r.Group("/api/networks")
	s.NewNetworkRouter(networkRouter)

	componentRouter := r.Group("/api/components")
	s.NewComponentRouter(componentRouter)

	r.Run(s.cfg.Host + ":" + s.cfg.Port)
}
