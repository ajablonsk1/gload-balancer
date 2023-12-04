package model

import (
	"sync/atomic"
)

type LoadDistributionStrategy interface {
	GetServer(serverPool *ServerPool) *Server
}

type RoundRobin struct{}

func (r *RoundRobin) GetServer(serverPool *ServerPool) *Server {
	next := serverPool.NextIndex()
	serversLength := len(serverPool.Servers)
	// getting adjusted length in order to loop through all servers
	fullCycleLength := serversLength + next
	for i := next; i < fullCycleLength; i++ {
		// normalizing index to be in slice range
		idx := i % serversLength
		server := serverPool.Servers[idx]
		if server.IsAlive() {
			// updating index if it was not the original one
			if i != next {
				atomic.StoreUint64(&serverPool.CurrentIdx, uint64(i))
			}
			return server
		}
	}
	return nil
}

type WeightedRoundRobin struct {
	sentReqToSameServer int
}

func (wR *WeightedRoundRobin) GetServer(serverPool *ServerPool) *Server {
	currServer := serverPool.Servers[serverPool.GetCurrentIdx()]
	// if current server is alive and we didn't send him enough requests, we return this server and add one to sent variable
	if wR.sentReqToSameServer < currServer.Weight && currServer.IsAlive() {
		wR.sentReqToSameServer = wR.sentReqToSameServer + 1
		return currServer
	}

	// if current server got enough requests or is dead we get next server
	// servers are sorted by weight so code will be simmilar to normal round robin
	wR.sentReqToSameServer = 0
	next := serverPool.NextIndex()
	serversLength := len(serverPool.Servers)
	fullCycleLength := serversLength + next
	for i := next; i < fullCycleLength; i++ {
		idx := i % serversLength
		server := serverPool.Servers[idx]
		if server.IsAlive() {
			if i != next {
				atomic.StoreUint64(&serverPool.CurrentIdx, uint64(i))
			}
			return server
		}
	}
	return nil
}

type IPHash struct{}

func (i *IPHash) GetServer(serverPool *ServerPool) *Server {
	// TODO
	return nil
}

type LeastConnection struct{}

func (l *LeastConnection) GetServer(serverPool *ServerPool) *Server {
	// TODO
	return nil
}

type WeightedLeastConnection struct{}

func (wL *WeightedLeastConnection) GetServer(serverPool *ServerPool) *Server {
	// TODO
	return nil
}

type WeightedResponseTime struct{}

func (l *WeightedResponseTime) GetServer(serverPool *ServerPool) *Server {
	// TODO
	return nil
}
