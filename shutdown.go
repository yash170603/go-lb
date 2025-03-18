package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// HandleShutdown manages graceful shutdown of the load balancer
func (lb *LoadBalancer) HandleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down load balancer gracefully...")
	os.Exit(0)
}
