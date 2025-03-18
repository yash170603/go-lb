package main

import (
	"log"
	"net/http"
	"sync"
)

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
	mutex           sync.Mutex
}

func NewLoadBalancer(port string) *LoadBalancer {
	lb := &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         make([]Server, 0),
	}

	go lb.healthCheckLoop()
	return lb
}

func (lb *LoadBalancer) AddServer(address string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	server := newSimpleServer(address)

	if server != nil {
		lb.servers = append(lb.servers, server)
		log.Printf("Added new Server : %s\n", address)
	}
}

func (lb *LoadBalancer) RemoveServer(address string) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for i, server := range lb.servers {
		if server.Address() == address {
			lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
			log.Printf("Removed Server : %s\n", address)
			return
		}
	}
}

func (lb *LoadBalancer) GetNextAvailableServer() Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.servers) == 0 {
		return nil
	}

	for range lb.servers {
		server := lb.servers[lb.roundRobinCount%len(lb.servers)]
		if server.IsAlive() {
			lb.roundRobinCount++
			return server
		}
		lb.roundRobinCount++
	}

	return nil
}

// ServeProxy forwards requests to the backend server
func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	server := lb.GetNextAvailableServer()
	if server == nil {
		http.Error(rw, "No available servers", http.StatusServiceUnavailable)
		return
	}

	log.Printf("Forwarding request to: %s\n", server.Address())
	server.Serve(rw, req)
}
