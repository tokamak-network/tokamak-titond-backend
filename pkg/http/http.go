package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tokamak-network/tokamak-titond-backend/api"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/api"
)

type Config struct {
	Host string
	Port string
}

type HTTPServer struct {
	R    *gin.Engine
	cfg  *Config
	apis *api.TitondAPI
}

func NewHTTPServer(cfg *Config, apis *api.TitondAPI) *HTTPServer {
	server := &HTTPServer{
		cfg:  cfg,
		apis: apis,
	}
	server.initialize()
	return server
}

func (s *HTTPServer) initialize() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	networkRouter := r.Group("/api/networks")
	s.NewNetworkRouter(networkRouter)

	componentRouter := r.Group("/api/components")
	s.NewComponentRouter(componentRouter)

	s.R = r
}

func (s *HTTPServer) Run() {

	s.R.Run(s.cfg.Host + ":" + s.cfg.Port)
}
