package model

import "sync/atomic"

type LoadDistributionStrategy interface {
	GetServer(serverPool *ServerPool) *Server
}

type RoundRobin struct{}

func (r *RoundRobin) GetServer(serverPool *ServerPool) *Server {
	next := serverPool.NextIndex()
	serversLength := len(serverPool.servers)
	// getting adjusted length in order to loop through all servers
	fullCycleLength := serversLength + next
	for i := next; i < fullCycleLength; i++ {
		// normalizing index to be in slice range
		idx := i % serversLength
		server := serverPool.servers[idx]
		if server.IsAlive() {
			// updating index if it was not the original one
			if i != next {
				atomic.StoreUint64(&serverPool.currentIdx, uint64(i))
			}
			return server
		}
	}
	return nil
}
