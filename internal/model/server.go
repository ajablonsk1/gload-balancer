package model

import (
	"cmp"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"go.uber.org/atomic"
	"time"
)

type Server struct {
	Url            *url.URL
	Alive          *atomic.Bool
	Proxy          *httputil.ReverseProxy
	Weight         int
	StickySessions map[string]time.Time
}

func (s *Server) IsAlive() bool {
	return s.Alive.Load()
}

func (s *Server) SetAlive(isAlive bool) {
	s.Alive.Store(isAlive)
}

func (s *Server) AddStickySession(remoteAddr string) {
	s.StickySessions[remoteAddr] = time.Now()
}

func (s *Server) HasStickySession(remoteAddr string) bool {
	_, ok := s.StickySessions[remoteAddr]
	return ok
}

func (s *Server) UpdateTimeForStickySession(remoteAddr string) {
	s.StickySessions[remoteAddr] = time.Now()
}

func (s *Server) DeleteStickySessionsIfTimeExpired() {
	for remoteAddr, stickyTime := range s.StickySessions {
		if stickyTime.Add(10 * time.Minute).Before(time.Now()) {
			delete(s.StickySessions, remoteAddr)
		}
	}
}

func (s *Server) checkHealth() {
    client := &http.Client{
        Timeout: 2 * time.Second,
    }

    resp, err := client.Get(s.Url.String())
    if err != nil {
        s.SetAlive(false)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        s.SetAlive(true)
    } else {
        s.SetAlive(false)
    }
}

func (s *Server) NumberOfStickySessions() int {
	return len(s.StickySessions)
}

func SortByWeight(a, b *Server) int {
	return cmp.Compare(b.Weight, a.Weight)
}

func SortByNConnections(a, b *Server) int {
	return cmp.Compare(a.NumberOfStickySessions(), b.NumberOfStickySessions())
}

type ServerPool struct {
	Servers    []*Server
	CurrentIdx atomic.Uint64
}

func (s *ServerPool) NextIndex() int {
	return int(s.CurrentIdx.Add(1)) % len(s.Servers)
}

func (s *ServerPool) GetCurrentIdx() int {
	return int(s.CurrentIdx.Load()) % len(s.Servers)
}

func (s *ServerPool) OrganizeStickySessions() {
	for _, server := range s.Servers {
		if !server.IsAlive() {
			server.StickySessions = make(map[string]time.Time)
			continue
		}
		server.DeleteStickySessionsIfTimeExpired()
	}
}

func (s *ServerPool) GetServerFromStickySession(remoteAddr string) *Server {
	s.OrganizeStickySessions()
	for _, server := range s.Servers {
		if server.HasStickySession(remoteAddr) {
			server.UpdateTimeForStickySession(remoteAddr)
			return server
		}
	}
	return nil
}

func (s *ServerPool) HealthCheck() {
	var wg sync.WaitGroup

	for _, currServer := range s.Servers {
		wg.Add(1)
		go func(server *Server) {
			defer wg.Done()
			server.checkHealth()
		}(currServer)
	}
	wg.Wait()
}
