package server

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func (s *BuffetAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s.Router.ServeHTTP(w, r)
	end := time.Now()
	log.WithFields(log.Fields{"Method": r.Method, "Path": r.URL, "Addr": r.RemoteAddr, "Elapsed": end.Sub(start)}).Info("")
}

func (s *BuffetAPIServer) InstallHandlers() {
	s.Router.HandleFunc("/", s.home)
	s.Router.HandleFunc("/ifttt/v1/actions/collect", s.collect)
	s.Router.HandleFunc("/ifttt/v1/test/setup", s.test)
	s.Router.HandleFunc("/ifttt/v1/status", s.status)
}
