package server

import (
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
}

func NewHttpServer() *HttpServer {
	srv := &HttpServer{
		Server: http.Server{
			Addr:         ":8080",
			ReadTimeout:  timeoutMax,
			WriteTimeout: timeoutMax,
		},
		identityService: identifier.NewIdentifyService(identifierMax),
	}
	return registerRequiredRoutes(srv)
}
