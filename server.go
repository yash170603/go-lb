package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
	HealthCheck()
}

type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
	alive   bool
	mutex   sync.RWMutex
}

func newSimpleServer(address string) *simpleServer {
	serverURL, err := url.Parse(address)
	if err != nil {
		log.Printf("Error parsing URL: %S\n", err)
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(serverURL)
	return &simpleServer{
		address: address,
		proxy:   proxy,
		alive:   true,
	}
}

func (s *simpleServer) Address() string {
	return s.address
}

func (s *simpleServer) IsAlive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.alive
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (s *simpleServer) HealthCheck() {
	resp, err := http.Get(s.address)
	s.mutex.Lock()
	if err != nil || resp.StatusCode >= 500 {
		log.Printf("Server %s is down\n", s.address)
		s.alive = false
	} else {
		s.alive = true
	}
	s.mutex.Unlock()
}
