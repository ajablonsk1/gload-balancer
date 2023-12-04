package model

import (
	"cmp"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Server struct {
	Url    *url.URL
	Alive  *atomic.Bool
	Proxy  *httputil.ReverseProxy
	Weight int
}

func (s *Server) IsAlive() bool {
	return s.Alive.Load()
}

func (s *Server) SetAlive(isAlive bool) {
	s.Alive.Swap(isAlive)
}

func SortByWeight(a, b *Server) int {
	return cmp.Compare(b.Weight, a.Weight)
}

type ServerPool struct {
	Servers    []*Server
	CurrentIdx uint64
}

func (s *ServerPool) GetAvailableServer(strategy LoadDistributionStrategy) *Server {
	return strategy.GetServer(s)
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.CurrentIdx, 1) % uint64(len(s.Servers)))
}

func (s *ServerPool) GetCurrentIdx() int {
	return int(s.CurrentIdx) % len(s.Servers)
}
