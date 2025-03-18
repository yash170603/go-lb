package main

import (
	"log"
	"net/http"
)

func main() {
	lb := NewLoadBalancer(LOAD_BALANCER_PORT)

	// Register API endpoints for server management
	http.HandleFunc("/register", lb.RegisterServerHandler)
	http.HandleFunc("/deregister", lb.DeregisterServerHandler)

	// Default request forwarding
	http.HandleFunc("/", lb.ServeProxy)

	// Handle graceful shutdown
	go lb.HandleShutdown()

	log.Printf("Load balancer is running on %s\n", LOAD_BALANCER_PORT)
	log.Fatal(http.ListenAndServe(LOAD_BALANCER_PORT, nil))
}
