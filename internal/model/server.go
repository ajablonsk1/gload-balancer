package model

import (
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Server struct {
	url          *url.URL
	alive        atomic.Bool
	reverseProxy *httputil.ReverseProxy
}

func (s *Server) IsAlive() bool {
	return s.alive.Load()
}

func (s *Server) SetAlive(isAlive bool) {
	s.alive.Swap(isAlive)
}

type ServerPool struct {
	servers    []*Server
	currentIdx uint64
}

func (s *ServerPool) GetAvailableServer(strategy LoadDistributionStrategy) *Server {
	return strategy.GetServer(s)
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.currentIdx, 1) % uint64(len(s.servers)))
}
