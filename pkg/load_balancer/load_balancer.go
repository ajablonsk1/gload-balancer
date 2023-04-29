package load_balancer

import (
	c "github.com/ajablonsk1/gload-balancer/internal/config"
	"github.com/ajablonsk1/gload-balancer/internal/model"
)

type LoadBalancer struct {
	strategy   model.LoadDistributionStrategy
	serverPool *model.ServerPool
}

func NewLoadBalancer(path string) (*LoadBalancer, error) {
	config, err := c.GetConfig(path)
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

	return &LoadBalancer{
		strategy:   strategy,
		serverPool: serverPool,
	}, nil
}

func (l *LoadBalancer) Start() {
	// TODO
}
