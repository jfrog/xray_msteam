package main

import (
	"net/http"
)

var routes map[string]func(http.ResponseWriter, *http.Request) = map[string]func(http.ResponseWriter, *http.Request){
	"/api/send":                 SendPage,
	"/api/ping":                 PingPage,
}

type serverHandler struct {
}

func (h *serverHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handler, found := routes[req.URL.Path]
	if !found {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte("404 page not found"))
		return
	}
	handler(res, req)
}

func GetWebServer() *http.Server {
	return &http.Server{Addr: ":8080", Handler: &serverHandler{}}
}
