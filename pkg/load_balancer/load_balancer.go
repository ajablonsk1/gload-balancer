package load_balancer

import (
	"log"
	"net/http"
	"net/url"

	c "github.com/ajablonsk1/gload-balancer/internal/config"
	"github.com/ajablonsk1/gload-balancer/internal/handler"
)

type LoadBalancer struct {
	addr         string
	proxyHandler http.Handler
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
		addr:         url.String(),
		proxyHandler: proxyHandler,
	}, nil
}

func (l *LoadBalancer) Start() {
	server := &http.Server{
		Addr:    l.addr,
		Handler: l.proxyHandler,
	}

	log.Fatal(server.ListenAndServe())
}
