package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}
type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

func newSimpleServer(address string) *simpleServer {
	serverUrl, err := url.Parse(address)
	proxyUrl := httputil.NewSingleHostReverseProxy(serverUrl)
	fmt.Println("Proxy URL: ", proxyUrl)
	handleErr(err)
	return &simpleServer{
		address: address,
		proxy:   proxyUrl, // this creates a reverse proxy for a particular p server

	}
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {

	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf(("error: %s\n"), err)
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string {
	return s.address
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (lb *LoadBalancer) getnextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {

	targetServer := lb.getnextAvailableServer()
	fmt.Printf(("forwarding request to address %s\n"), targetServer.Address())
	targetServer.Serve(rw, r)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.google.com/"),
		newSimpleServer("http://www.instagram.com"),
		newSimpleServer("https://www.youtube.com/"),
	}

	lb := NewLoadBalancer(":8000", servers)
	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Load balancer listening on %s\n", lb.port)
	http.ListenAndServe(lb.port, nil)
}
