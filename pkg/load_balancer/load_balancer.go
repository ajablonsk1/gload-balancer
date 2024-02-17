package load_balancer

import (
	"log"
	"net/http"
	"net/url"
	"time"

	c "github.com/ajablonsk1/gload-balancer/internal/config"
	"github.com/ajablonsk1/gload-balancer/internal/handler"
)

type LoadBalancer struct {
	Addr         string
	ProxyHandler *handler.ProxyHandler
}

func NewLoadBalancer(path string) (*LoadBalancer, error) {
	config, err := c.GetConfig(path)
	if err != nil {
		return nil, err
	}

	addr, err := config.GetAddress()
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	strategy, err := config.GetLoadStrategy()
	if err != nil {
		return nil, err
	}

	serverPool, err := config.GetServerPool()
	if err != nil {
		return nil, err
	}

	proxyHandler := &handler.ProxyHandler{
		Strategy:   strategy,
		ServerPool: serverPool,
	}

	return &LoadBalancer{
		Addr:         url.String(),
		ProxyHandler: proxyHandler,
	}, nil
}

func (l *LoadBalancer) RunHealthChecks() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		l.ProxyHandler.ServerPool.HealthCheck()
	}
}

func (l *LoadBalancer) Start() {
	server := &http.Server{
		Addr:    l.Addr,
		Handler: l.ProxyHandler,
	}

	go l.RunHealthChecks()

	log.Fatal(server.ListenAndServe())
}
