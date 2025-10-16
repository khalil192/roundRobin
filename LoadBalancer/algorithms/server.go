package algorithms

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"
	"sync/atomic"
)

type Server struct {
	URL     *url.URL
	healthy bool
	mu      sync.RWMutex
	//tracking metrics for tests
	servedYet uint32
	failedYet uint32
}

type ServerList []*Server

func NewServer(hostURL *url.URL) *Server {
	return &Server{
		URL:       hostURL,
		healthy:   true,
		mu:        sync.RWMutex{},
		servedYet: 0,
		failedYet: 0,
	}
}

func NewServerList(serverUrls []string) []*Server {
	servers := make([]*Server, 0, len(serverUrls))

	for _, rawUrl := range serverUrls {
		parsedUrl, parseErr := url.Parse(rawUrl)
		if parseErr != nil {
			log.Printf("Failed to parse URL: %s, error:%s", rawUrl, parseErr.Error())
			continue
		}
		servers = append(servers, NewServer(parsedUrl))
	}

	return servers
}

func (s *Server) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.healthy
}

func (s *Server) HandleReqServed() {
	atomic.AddUint32(&s.servedYet, 1)
}

func (s *Server) HandleReqFailed() {
	atomic.AddUint32(&s.failedYet, 1)
}

func (s *Server) setHealthy(status bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.healthy = status
}

type contextKey string

const serverContextKey = contextKey("BackendServer")

// NewServerContext add the server pointer to the request context
func NewServerContext(ctx context.Context, server *Server) context.Context {
	return context.WithValue(ctx, serverContextKey, server)
}

// GetServerFromContext retrieve the stored server from the request context
func GetServerFromContext(ctx context.Context) *Server {
	if srv, ok := ctx.Value(serverContextKey).(*Server); ok {
		return srv
	}
	return nil
}

type Stats struct {
	URL      string
	Healthy  string
	Requests string
	Failed   string
}

func (s *Server) GetWritableStats() []string {
	return []string{
		s.URL.String(),
		fmt.Sprint(s.IsHealthy()),
		fmt.Sprint(atomic.LoadUint32(&s.servedYet)),
		fmt.Sprint(atomic.LoadUint32(&s.failedYet)),
	}
}

func (sl ServerList) FindServerIndex(target *Server) int {
	for i, server := range sl {
		if server == target {
			return i
		}
	}

	return -1
}
