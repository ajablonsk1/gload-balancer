package model

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"go.uber.org/atomic"
)

func TestServerIsAlive(t *testing.T) {
	server := &Server{
		Alive: atomic.NewBool(true),
	}

	if !server.IsAlive() {
		t.Error("Expected server to be alive")
	}
}

func TestServerSetAlive(t *testing.T) {
	isAlive := &atomic.Bool{}
	isAlive.Store(true) 

	server := &Server{
		Alive: atomic.NewBool(true),
	}

	server.SetAlive(false)

	if server.IsAlive() {
		t.Error("Expected server to be dead")
	}
}

func TestServerAddStickySession(t *testing.T) {
	server := &Server{
		StickySessions: make(map[string]time.Time),
	}

	remoteAddr := "127.0.0.1"
	server.AddStickySession(remoteAddr)

	if _, ok := server.StickySessions[remoteAddr]; !ok {
		t.Errorf("Expected sticky session for remote address %s", remoteAddr)
	}
}

func TestServerHasStickySession(t *testing.T) {
	server := &Server{
		StickySessions: make(map[string]time.Time),
	}

	remoteAddr := "127.0.0.1"
	server.StickySessions[remoteAddr] = time.Now()

	if !server.HasStickySession(remoteAddr) {
		t.Errorf("Expected server to have sticky session for remote address %s", remoteAddr)
	}
}

func TestServerUpdateTimeForStickySession(t *testing.T) {
	server := &Server{
		StickySessions: make(map[string]time.Time),
	}

	remoteAddr := "127.0.0.1"
	server.StickySessions[remoteAddr] = time.Now()

	server.UpdateTimeForStickySession(remoteAddr)

	if server.StickySessions[remoteAddr].Equal(time.Now()) {
		t.Errorf("Expected sticky session time to be updated for remote address %s", remoteAddr)
	}
}

func TestServerDeleteStickySessionsIfTimeExpired(t *testing.T) {
	server := &Server{
		StickySessions: make(map[string]time.Time),
	}

	remoteAddr := "127.0.0.1"
	server.StickySessions[remoteAddr] = time.Now().Add(-time.Hour)

	server.DeleteStickySessionsIfTimeExpired()

	if _, ok := server.StickySessions[remoteAddr]; ok {
		t.Errorf("Expected sticky session for remote address %s to be deleted", remoteAddr)
	}
}

func TestServerNumberOfStickySessions(t *testing.T) {
	server := &Server{
		StickySessions: make(map[string]time.Time),
	}

	server.StickySessions["127.0.0.1"] = time.Now()
	server.StickySessions["192.168.0.1"] = time.Now()

	count := server.NumberOfStickySessions()

	if count != 2 {
		t.Errorf("Expected number of sticky sessions to be 2, got %d", count)
	}
}

func TestServerPoolNextIndex(t *testing.T) {
	serverPool := &ServerPool{
		Servers: []*Server{
			{Weight: 1},
			{Weight: 2},
			{Weight: 3},
		},
		CurrentIdx: *atomic.NewUint64(0),
	}

	index := serverPool.NextIndex()

	if index != 1 {
		t.Errorf("Expected next index to be 1, got %d", index)
	}
}

func TestServerPoolGetCurrentIdx(t *testing.T) {
	serverPool := &ServerPool{
		CurrentIdx: *atomic.NewUint64(2),
		Servers:    []*Server{{}, {}, {}},
	}

	index := serverPool.GetCurrentIdx()

	if index != 2 {
		t.Errorf("Expected current index to be 2, got %d", index)
	}
}

func TestServerPoolOrganizeStickySessions(t *testing.T) {
	alive := &atomic.Bool{}
	alive.Store(true)

	dead := &atomic.Bool{}
	dead.Store(false)

	serverPool := &ServerPool{
		Servers: []*Server{
			{StickySessions: map[string]time.Time{"127.0.0.1": time.Now()}, Alive: atomic.NewBool(false)},
			{StickySessions: map[string]time.Time{"192.168.0.1": time.Now()}, Alive: atomic.NewBool(true)},
		},
	}
	serverPool.Servers[0].SetAlive(false)
	serverPool.Servers[1].SetAlive(true)

	serverPool.OrganizeStickySessions()

	if len(serverPool.Servers[0].StickySessions) != 0 {
		t.Error("Expected first server to have no sticky sessions")
	}

	if len(serverPool.Servers[1].StickySessions) != 1 {
		t.Error("Expected second server to have one sticky session")
	}
}

func TestServerPoolGetServerFromStickySession(t *testing.T) {
	alive := &atomic.Bool{}
	alive.Store(true)

	serverPool := &ServerPool{
		Servers: []*Server{
			{StickySessions: map[string]time.Time{"127.0.0.1": time.Now()}, Alive: atomic.NewBool(true)},
			{StickySessions: map[string]time.Time{"192.168.0.1": time.Now()}, Alive: atomic.NewBool(true)},
		},
	}

	remoteAddr := "192.168.0.1"
	server := serverPool.GetServerFromStickySession(remoteAddr)

	if server == nil {
		t.Errorf("Expected server to be found for remote address %s", remoteAddr)
	}
}

func TestServerPoolHealthCheck(t *testing.T) {
	url1, _ := url.Parse("http://127.0.0.1:8080")
	url2, _ := url.Parse("http://127.0.0.1:8081")

	serverPool := &ServerPool{
		Servers: []*Server{
			{Url: url1, Alive: atomic.NewBool(true)},
			{Url: url2, Alive: atomic.NewBool(true)},
		},
	}

	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    	w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()


	sUrl, _ := url.Parse(testServer.URL)
	serverPool.Servers[0].Url = sUrl

	serverPool.HealthCheck()

	if !serverPool.Servers[0].IsAlive() {
		t.Error("Expected first server to be alive")
	}

	if serverPool.Servers[1].IsAlive() {
		t.Error("Expected second server to be dead")
	}
}
