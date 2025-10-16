package algorithms

import "sync/atomic"

type RoundRobinStruct struct {
	servers []*Server
	current uint32
}

func NewAtomicRoundRobinBalancer(servers []*Server) *RoundRobinStruct {
	return &RoundRobinStruct{
		servers: servers,
		current: 0,
	}
}

func (lb *RoundRobinStruct) GetNextHealthyServer() *Server {
	for range lb.servers {
		// lb.current is atomically incremented and no two concurrent req will arrive at same nextIdx
		nextIdx := int(atomic.AddUint32(&lb.current, 1))
		idxOnList := (nextIdx) % (len(lb.servers))
		if lb.servers[idxOnList].IsHealthy() {
			return lb.servers[idxOnList]
		}
	}

	//no healthy server
	return nil
}

func (lb *RoundRobinStruct) GetAllServers() []*Server {
	return lb.servers
}

func (lb *RoundRobinStruct) SetServerStatus(server *Server, isHealthy bool) {
	server.setHealthy(isHealthy)
}

func (lb *RoundRobinStruct) CompleteHeathCheck() {
}
