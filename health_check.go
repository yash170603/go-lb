package main

import (
	"log"
	"time"
)

// HealthCheckLoop periodically checks the health of servers
func (lb *LoadBalancer) healthCheckLoop() {
	ticker := time.NewTicker(time.Duration(HEALTH_CHECK_INTERVAL) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running health checks...")
		lb.mutex.Lock()
		for _, server := range lb.servers {
			server.HealthCheck()
		}
		lb.mutex.Unlock()
	}
}
