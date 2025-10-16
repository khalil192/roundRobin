package algorithms

import (
	"sync"
)

type GlobalLockRoundRobinStruct struct {
	servers []*Server
	current uint32
	mutex   sync.Mutex
}

func NewGlobalLockRoundRobin(servers []*Server) *GlobalLockRoundRobinStruct {
	return &GlobalLockRoundRobinStruct{
		servers: servers,
		current: 0,
	}
}

func (lb *GlobalLockRoundRobinStruct) GetNextHealthyServer() *Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	for range lb.servers {
		lb.current += 1
		idxOnList := (int(lb.current)) % (len(lb.servers))
		if lb.servers[idxOnList].IsHealthy() {
			return lb.servers[idxOnList]
		}
	}

	//no healthy server
	return nil
}

func (lb *GlobalLockRoundRobinStruct) SetServerStatus(server *Server, isHealthy bool) {
	server.setHealthy(isHealthy)
}

func (lb *GlobalLockRoundRobinStruct) GetAllServers() []*Server {
	return lb.servers
}

func (lb *GlobalLockRoundRobinStruct) CompleteHeathCheck() {
}
