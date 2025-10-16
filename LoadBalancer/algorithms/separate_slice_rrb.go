package algorithms

import (
	"sync"
	"sync/atomic"
)

type HealthySlice struct {
	mu               sync.RWMutex
	healthyServers   ServerList
	unhealthyServers ServerList
	current          uint64
}

func (lb *HealthySlice) CompleteHeathCheck() {
}

func NewSeparateSlice(serverList []*Server) *HealthySlice {
	return &HealthySlice{
		healthyServers:   serverList,
		unhealthyServers: make([]*Server, 0),
		current:          0,
	}
}

func (lb *HealthySlice) GetNextHealthyServer() *Server {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	numHealthy := len(lb.healthyServers)
	if numHealthy == 0 {
		return nil
	}

	nextIndex := atomic.AddUint64(&lb.current, 1)

	idx := int((nextIndex - 1) % uint64(numHealthy))
	return lb.healthyServers[idx]
}

func (lb *HealthySlice) SetServerStatus(server *Server, isHealthy bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if isHealthy && !server.IsHealthy() {
		lb.unhealthyServers, lb.healthyServers = lb.transferElementFrom(lb.unhealthyServers, lb.healthyServers, server)
		server.setHealthy(true)
	} else if !isHealthy && server.IsHealthy() {
		lb.healthyServers, lb.unhealthyServers = lb.transferElementFrom(lb.healthyServers, lb.unhealthyServers, server)
		server.setHealthy(false)
	}

}

func (lb *HealthySlice) transferElementFrom(src ServerList, dest ServerList, target *Server) (ServerList, ServerList) {
	idx := src.FindServerIndex(target)
	if idx != -1 {
		dest = append(dest, target)
		src[idx] = src[len(src)-1]
		src = src[:len(src)-1]
	}

	return src, dest
}

func (lb *HealthySlice) GetAllServers() []*Server {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	allServers := make([]*Server, 0, len(lb.healthyServers)+len(lb.unhealthyServers))
	allServers = append(allServers, lb.healthyServers...)
	allServers = append(allServers, lb.unhealthyServers...)
	return allServers
}
