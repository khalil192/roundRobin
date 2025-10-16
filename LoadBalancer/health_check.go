package main

import (
	"log"
	"net/http"
	"roundrobin/LoadBalancer/algorithms"
	"time"
)

var healthCheckClient = &http.Client{
	Timeout: 3 * time.Second,
}

func StartHealthChecks(lb LoadBalancerI, interval time.Duration) {
	go func() {
		log.Printf("health checker started with a %s interval...", interval)

		ticker := time.NewTicker(interval)

		for range ticker.C {
			servers := lb.GetAllServers()
			//log.Printf("checking %d servers health", len(servers))
			for _, server := range servers {
				//log.Println("checking server:", server.URL.String())
				checkServerHealth(lb, server)
			}
		}

		lb.CompleteHeathCheck()
	}()
}

func checkServerHealth(lb LoadBalancerI, server *algorithms.Server) {
	// The health check endpoint is assumed to be at the "/health" path.
	healthURL := server.URL.String() + "/health"

	resp, err := healthCheckClient.Get(healthURL)
	currentStatusIsHealthy := server.IsHealthy()

	if err != nil {
		if currentStatusIsHealthy {
			log.Printf("Health check failed for %s. It is now unhealthy.", server.URL)
			lb.SetServerStatus(server, false)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if !currentStatusIsHealthy {
			log.Printf("Server %s is marked as healthy after health check", server.URL)
			lb.SetServerStatus(server, true)
		}
	} else {
		if currentStatusIsHealthy {
			log.Printf("Server %s is marked as unhealthy after health check, response code: %d", server.URL, resp.StatusCode)
			lb.SetServerStatus(server, false)
		}
	}
}
