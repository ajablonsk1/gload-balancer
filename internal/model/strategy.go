package model

import (
	"slices"
	"sync/atomic"

	"github.com/ajablonsk1/gload-balancer/internal/utils"
)

type LoadDistributionStrategy interface {
	GetServer(serverPool *ServerPool, remoteAddr string) *Server
}

type RoundRobin struct{}

func (r *RoundRobin) GetServer(serverPool *ServerPool, remoteAddr string) *Server {
	if s := serverPool.GetServerFromStickySession(remoteAddr); s != nil {
		return s
	}

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
			server.AddStickySession(remoteAddr)
			return server
		}
	}
	return nil
}

type WeightedRoundRobin struct {
	sentReqToSameServer int
}

func (wR *WeightedRoundRobin) GetServer(serverPool *ServerPool, remoteAddr string) *Server {
	if s := serverPool.GetServerFromStickySession(remoteAddr); s != nil {
		return s
	}

	currServer := serverPool.Servers[serverPool.GetCurrentIdx()]
	// if current server is alive and we didn't send him enough requests, we return this server and add one to sent variable
	if wR.sentReqToSameServer < currServer.Weight && currServer.IsAlive() {
		wR.sentReqToSameServer = wR.sentReqToSameServer + 1
		currServer.AddStickySession(remoteAddr)
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
			server.AddStickySession(remoteAddr)
			return server
		}
	}
	return nil
}

type IPHash struct{}

func (i *IPHash) GetServer(serverPool *ServerPool, remoteAddr string) *Server {
	if s := serverPool.GetServerFromStickySession(remoteAddr); s != nil {
		return s
	}

	hash := utils.Hash(remoteAddr)
	serversLength := len(serverPool.Servers)
	// we get idx of server based on hash
	idx := int(hash) % serversLength

	// we do a full cycle to find alive server, if server based on hash is dead we want to find next alive server
	fullCycleLength := serversLength + idx
	for i := idx; i < fullCycleLength; i++ {
		idx := i % serversLength
		server := serverPool.Servers[idx]
		if server.IsAlive() {
			server.AddStickySession(remoteAddr)
			return server
		}
	}
	return nil
}

type LeastSession struct{}

func (l *LeastSession) GetServer(serverPool *ServerPool, remoteAddr string) *Server {
	if s := serverPool.GetServerFromStickySession(remoteAddr); s != nil {
		return s
	}

	slices.SortFunc(serverPool.Servers, SortByNConnections)
	serversLength := len(serverPool.Servers)
	for i := 0; i < serversLength; i++ {
		server := serverPool.Servers[i]
		if server.IsAlive() {
			server.AddStickySession(remoteAddr)
			return server
		}
	}
	return nil
}

type WeightedLeastSession struct{}

func (wL *WeightedLeastSession) GetServer(serverPool *ServerPool, remoteAddr string) *Server {
	if s := serverPool.GetServerFromStickySession(remoteAddr); s != nil {
		return s
	}

	slices.SortFunc(serverPool.Servers, SortByNConnections)
	serversLength := len(serverPool.Servers)
	serversGroupedByNConnections := make([][]*Server, 0)
	lastNConnections := -1
	for i := 0; i < serversLength; i++ {
		server := serverPool.Servers[i]
		if server.IsAlive() {
			if lastNConnections != server.NumberOfStickySessions() {
				serversGroupedByNConnections = append(serversGroupedByNConnections, make([]*Server, 0))
				lastNConnections = server.NumberOfStickySessions()
			}
			serversGroupedByNConnections[len(serversGroupedByNConnections)-1] = append(serversGroupedByNConnections[len(serversGroupedByNConnections)-1], server)
		}
	}
	for _, servers := range serversGroupedByNConnections {
		if len(servers) > 0 {
			slices.SortFunc(servers, SortByWeight)
			server := servers[0]
			server.AddStickySession(remoteAddr)
			return server
		}
	}
	return nil
}

type WeightedResponseTime struct{}

func (l *WeightedResponseTime) GetServer(serverPool *ServerPool, remoteAddr string) *Server {
	// TODO
	return nil
}
