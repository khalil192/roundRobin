package algorithms

import (
	"github.com/gammazero/deque"
	"sync"
)

type HealthyQueue struct {
	healthyServers deque.Deque[*Server]
	deadServers    []*Server
	mutex          sync.Mutex
}

func NewQueueRRB(servers []*Server) *HealthyQueue {
	sp := &HealthyQueue{
		deadServers:    make([]*Server, 0),
		healthyServers: deque.Deque[*Server]{},
	}

	for _, server := range servers {
		sp.healthyServers.PushBack(server)
	}
	return sp
}

func (lb *HealthyQueue) GetNextHealthyServer() *Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for range lb.healthyServers.Len() {
		if !lb.healthyServers.Front().healthy {
			lb.deadServers = append(lb.deadServers, lb.healthyServers.PopFront())
			continue
		}
		svr := lb.healthyServers.PopFront()
		lb.healthyServers.PushBack(svr)
		return svr
	}

	return nil
}

func (lb *HealthyQueue) GetAllServers() []*Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	servers := make([]*Server, 0)

	for _, server := range lb.deadServers {
		servers = append(servers, server)
	}

	for i := range lb.healthyServers.Len() {
		servers = append(servers, lb.healthyServers.At(i))
	}

	return servers
}

func (lb *HealthyQueue) SetServerStatus(server *Server, isHealthy bool) {
	server.setHealthy(isHealthy)
}

func (lb *HealthyQueue) CompleteHeathCheck() {
	updatedDeadServers := make([]*Server, 0)
	for _, server := range lb.deadServers {
		if server.healthy {
			lb.healthyServers.PushBack(server)
		} else {
			updatedDeadServers = append(updatedDeadServers, server)
		}
	}

	lb.deadServers = updatedDeadServers
}
