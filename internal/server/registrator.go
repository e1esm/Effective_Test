package server

import "net/http"

func registerRequiredRoutes(srv *HttpServer) *HttpServer {
	r := http.NewServeMux()
	r.HandleFunc("/api/add", srv.New)
	r.HandleFunc("/api/delete", srv.Delete)
	r.HandleFunc("/api/update", srv.Change)
	r.HandleFunc("/api/get", srv.Get)
	srv.Handler = r
	return srv
}
