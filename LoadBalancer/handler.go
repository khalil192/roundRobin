package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"roundrobin/LoadBalancer/algorithms"
)

type LoadBalancerI interface {
	GetNextHealthyServer() *algorithms.Server
	GetAllServers() []*algorithms.Server
	CompleteHeathCheck()
	SetServerStatus(server *algorithms.Server, isHealthy bool)
}

// NewProxyHandler creates a robust, reusable HTTP handler that acts as a reverse proxy.
func NewProxyHandler(lb LoadBalancerI) http.Handler {
	proxy := &httputil.ReverseProxy{}
	proxy.Director = NewDirectorFunc(lb)
	proxy.ErrorHandler = NewErrorHandlerFunc(lb)
	return proxy
}

func NewDirectorFunc(lb LoadBalancerI) func(req *http.Request) {
	return func(req *http.Request) {
		backendServer := lb.GetNextHealthyServer()
		if backendServer == nil {
			return
		}

		//log.Printf("Forwarding request to backend: %s", backendServer.URL)

		req.URL.Scheme = backendServer.URL.Scheme
		req.URL.Host = backendServer.URL.Host
		req.URL.Path = req.URL.Path
		req.Host = backendServer.URL.Host

		backendServer.HandleReqServed()
		*req = *req.WithContext(algorithms.NewServerContext(req.Context(), backendServer))
	}
}

func NewErrorHandlerFunc(lb LoadBalancerI) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		//log.Printf("Reverse Proxy error: %v", err)

		backendServer := algorithms.GetServerFromContext(r.Context())
		if backendServer != nil {
			lb.SetServerStatus(backendServer, false)
			backendServer.HandleReqFailed()
			//log.Printf("Marked server %s as unhealthy due to error.", backendServer.URL)
		}

		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		}
	}
}
