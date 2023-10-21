package server

import "net/http"

func registerRequiredRoutes(srv *HttpServer) *HttpServer {
	r := http.NewServeMux()
	r.HandleFunc("/api/add", srv.New)
	srv.Handler = r
	return srv
}
