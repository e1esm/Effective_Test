package server

import (
	"github.com/e1esm/Effective_Test/internal/repository/postgres"
	"github.com/e1esm/Effective_Test/internal/service/aggregator"
	"github.com/e1esm/Effective_Test/internal/service/identifier"
	"net/http"
	"time"
)

const (
	timeoutMax    = time.Second * 5
	identifierMax = time.Second * 2
)

type HttpServer struct {
	http.Server
	identityService identifier.Identifier
	userService     aggregator.Aggregator
}

func NewHttpServer() *HttpServer {
	srv := &HttpServer{
		Server: http.Server{
			Addr:         ":8080",
			ReadTimeout:  timeoutMax,
			WriteTimeout: timeoutMax,
		},
		identityService: identifier.NewIdentifyService(identifierMax),
		userService:     aggregator.NewUserService(postgres.NewPeopleRepository()),
	}
	return registerRequiredRoutes(srv)
}
