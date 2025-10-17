package algorithms_test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"roundrobin/LoadBalancer/algorithms"
	"testing"
	"time"
)

type LoadBalancerI interface {
	GetNextHealthyServer() *algorithms.Server
}

type constructorFunc func(servers algorithms.ServerList) LoadBalancerI

// basic unit-tests to check the correctness of the implementations
func TestRoundRobinStruct_GetNextHealthyServer(t *testing.T) {
	constructors := []constructorFunc{
		func(servers algorithms.ServerList) LoadBalancerI {
			return algorithms.NewGlobalLockRoundRobin(servers)
		},
		func(servers algorithms.ServerList) LoadBalancerI {
			return algorithms.NewSeparateSlice(servers)
		},
		func(servers algorithms.ServerList) LoadBalancerI {
			return algorithms.NewAtomicRoundRobinBalancer(servers)
		},
		func(servers algorithms.ServerList) LoadBalancerI {
			return algorithms.NewQueueRRB(servers)
		},
	}

	for _, constructor := range constructors {
		urls := []string{"localhost:8080", "localhost:8081", "localhost:8082"}
		t.Run("should assign servers each once", func(t *testing.T) {
			serverList := algorithms.NewServerList(urls)
			lb := constructor(serverList)

			s1 := lb.GetNextHealthyServer()
			time.Sleep(time.Millisecond * 100)
			s2 := lb.GetNextHealthyServer()
			time.Sleep(time.Millisecond * 100)
			s3 := lb.GetNextHealthyServer()

			assert.NotEqual(t, s1.URL.String(), s2.URL.String())
			assert.NotEqual(t, s1.URL.String(), s3.URL.String())
			assert.NotEqual(t, s2.URL.String(), s3.URL.String())

		})

		t.Run("should return nil if no server is healthy", func(t *testing.T) {
			serverList := algorithms.NewServerList(urls)
			lb := algorithms.NewAtomicRoundRobinBalancer(serverList)
			for i := range len(serverList) {
				lb.SetServerStatus(serverList[i], false)
			}

			server := lb.GetNextHealthyServer()
			assert.Equal(t, nil, server)

		})

		t.Run("should return only the healthy server", func(t *testing.T) {
			serverList := algorithms.NewServerList(urls)
			lb := algorithms.NewAtomicRoundRobinBalancer(serverList)
			lb.SetServerStatus(serverList[1], false)
			lb.SetServerStatus(serverList[2], false)

			s1 := lb.GetNextHealthyServer()
			time.Sleep(time.Millisecond * 100)
			s2 := lb.GetNextHealthyServer()
			time.Sleep(time.Millisecond * 100)
			s3 := lb.GetNextHealthyServer()

			assert.Equal(t, s1.URL.String(), urls[0])
			assert.Equal(t, s2.URL.String(), urls[0])
			assert.Equal(t, s3.URL.String(), urls[0])
		})

		t.Run("should return all servers irrespective of their health status", func(t *testing.T) {
			urls := make([]string, 0)
			for i := 0; i < 10; i++ {
				urls = append(urls, fmt.Sprintf("localhost:%d", 8080+i))
			}

			serverList := algorithms.NewServerList(urls)
			lb := algorithms.NewAtomicRoundRobinBalancer(serverList)

			for i := range serverList {
				if i%2 == 0 {
					lb.SetServerStatus(serverList[i], false)
				}
			}

			allServers := lb.GetAllServers()
			assert.Equal(t, len(serverList), len(allServers))
		})
	}

}
