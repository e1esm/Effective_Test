package server

import "net/http"

func registerRequiredRoutes(srv *HttpServer) *HttpServer {
	r := http.NewServeMux()
	r.HandleFunc("/api/add", srv.New)
	r.HandleFunc("/api/delete", srv.Delete)
	r.HandleFunc("/api/update", srv.Change)
	srv.Handler = r
	return srv
}
