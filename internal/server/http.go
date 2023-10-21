package server

import (
	"errors"
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

var (
	methodErr      = errors.New("invalid method for the URL: %v")
	invalidReq     = errors.New("invalid request was sent to the URL: %v")
	marshallingErr = errors.New("error while operating over the request input data")
	identityErr    = errors.New("error occurred while identifying the user: %v")
	saveErr        = errors.New("error while inserting user: %v")
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
