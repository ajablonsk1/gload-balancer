package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
	"slices"
	"go.uber.org/atomic"
	"time"

	"github.com/ajablonsk1/gload-balancer/internal/model"
)

type Config map[string]interface{}

func GetConfig(path string) (Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func (c Config) GetAddress() (string, error) {
	if address, ok := c["address"].(string); ok {
		return address, nil
	} else {
		return address, errors.New("wrong config or no address attribute provided")
	}
}

func (c Config) GetLoadStrategy() (model.LoadDistributionStrategy, error) {
	if loadStrategy, ok := c["strategy"].(string); ok {
		switch loadStrategy {
		case "round-robin":
			return &model.RoundRobin{}, nil
		case "weighted-round-robin":
			return &model.WeightedRoundRobin{}, nil
		case "ip-hash":
			return &model.IPHash{}, nil
		case "least-connection":
			return &model.LeastSession{}, nil
		case "weighted-least-connection":
			return &model.WeightedLeastSession{}, nil
		default:
			return nil, errors.New("wrong strategy type in config file")
		}
	} else {
		return nil, errors.New("wrong config or no strategy attribute provided")
	}
}

func (c Config) GetServerPool() (*model.ServerPool, error) {
	if serversJson, ok := c["servers"].([]interface{}); ok {
		if len(serversJson) < 1 {
			return nil, errors.New("there must be at least one server provided")
		}

		servers := make([]*model.Server, 0)
		for _, server := range serversJson {
			server, ok := server.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("error with type assertion. expected map[string]interface{} got %T", server)
			}

			serverUrl, err := c.getServerUrl(server)
			if err != nil {
				return nil, err
			}

			proxy := httputil.NewSingleHostReverseProxy(serverUrl)

			weight := c.getServerWeight(server)

			servers = append(servers, &model.Server{
				Url:            serverUrl,
				Alive:          atomic.NewBool(true),
				Proxy:          proxy,
				Weight:         weight,
				StickySessions: make(map[string]time.Time),
			})
		}

		strategy, err := c.GetLoadStrategy()
		if err != nil {
			return nil, err
		}

		if _, ok := strategy.(*model.WeightedRoundRobin); ok {
			slices.SortFunc(servers, model.SortByWeight)
		}
		
		return &model.ServerPool{
			Servers: servers,
		}, nil
	} else {
		return nil, errors.New("wrong config or no servers attribute provided")
	}
}

func (c Config) getServerUrl(server map[string]interface{}) (*url.URL, error) {
	host, ok := server["host"].(string)
	if !ok {
		return nil, fmt.Errorf("error with type assertion. expected string under 'host' key")
	}

	serverUrl, err := url.Parse("http://" + host)
	if err != nil {
		return nil, err
	}

	return serverUrl, nil
}

func (c Config) getServerWeight(server map[string]interface{}) int {
	weight, _ := server["weight"].(float64)
	return int(weight)
}
