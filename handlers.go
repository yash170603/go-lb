package main

import (
	"fmt"
	"net/http"
)

// RegisterServerHandler adds a new backend server
func (lb *LoadBalancer) RegisterServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	lb.AddServer(address)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server %s registered successfully\n", address)
}

// DeregisterServerHandler removes a backend server
func (lb *LoadBalancer) DeregisterServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	lb.RemoveServer(address)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server %s deregistered successfully\n", address)
}
