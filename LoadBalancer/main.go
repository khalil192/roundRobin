package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"roundrobin/LoadBalancer/algorithms"
	"syscall"
	"time"
)

func main() {

	port := flag.String("port", "9000", "Port to listen on")
	configFile := flag.String("config", "././all_servers.txt", "Path to config file")
	algoToUse := flag.String("algoToUse", "queue", "Path to config file")
	flag.Parse()

	lb := loadLoadBalancer(*configFile, *algoToUse)
	router := NewRouter(lb)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: router,
	}

	// Run server in a goroutine to support graceful shutdown
	go func() {
		log.Println("Server started on : ", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	StartHealthChecks(lb, time.Second*10)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown due to an error: %v", err)
	}

	logServerStatsToFile(lb.GetAllServers(), "server_stats.csv")
}

func getLoadBalancer(algorithm string, serverList []*algorithms.Server) LoadBalancerI {
	switch algorithm {
	case "lock":
		return algorithms.NewGlobalLockRoundRobin(serverList)
	case "atomic":
		return algorithms.NewAtomicRoundRobinBalancer(serverList)
	case "separate_slice":
		return algorithms.NewSeparateSlice(serverList)
	}
	return nil
}

func loadLoadBalancer(configFile string, algoToUse string) LoadBalancerI {
	cfg := Config(configFile)
	if len(cfg.BackendServerURLs) == 0 {
		log.Fatal("No backend servers configured.")
	}

	serverList := algorithms.NewServerList(cfg.BackendServerURLs)
	if len(serverList) == 0 {
		log.Fatal("Server list is empty.")
	}

	lb := getLoadBalancer(algoToUse, serverList)
	if lb == nil {
		panic("LoadBalancer Algorithm is not configured")
	}

	return lb
}

func NewRouter(lb LoadBalancerI) http.Handler {
	mux := http.NewServeMux()
	proxyHandler := NewProxyHandler(lb)
	mux.Handle("/", proxyHandler)

	return mux
}

func logServerStatsToFile(servers []*algorithms.Server, filename string) {
	log.Printf("Logging server stats to %s...", filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("ERROR: Could not create stats file: %v", err)
		return
	}
	defer file.Close()

	statsList := make([][]string, len(servers))
	for i, s := range servers {
		statsList[i] = s.GetWritableStats()
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"URL", "Healthy", "Served Requests", "Failed Requests"})

	for _, stats := range statsList {
		if err := writer.Write(stats); err != nil {
			log.Println("Error writing csv row:", err)
		}
	}

	log.Printf("Successfully wrote stats for %d servers.", len(servers))
}
