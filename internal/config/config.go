package config

import (
	"encoding/json"
	"errors"
	"github.com/ajablonsk1/gload-balancer/internal/model"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
)

type Config map[string]interface{}

func GetConfig(path string) (*Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
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
			return &model.LeastConnection{}, nil
		case "weighted-least-connection":
			return &model.WeightedLeastConnection{}, nil
		case "weighted-response-time":
			return &model.WeightedResponseTime{}, nil
		default:
			return nil, errors.New("wrong strategy type in config file")
		}
	} else {
		return nil, errors.New("no strategy attribute provided")
	}
}

func (c Config) GetServerPool() (*model.ServerPool, error) {
	if serversJson, ok := c["servers"].([]map[string]interface{}); ok {
		if len(serversJson) < 1 {
			return nil, errors.New("there must be at least one server provided")
		}

		servers := make([]*model.Server, 1)
		for _, server := range serversJson {
			serverUrl, err := url.Parse(server["host"].(string))
			if err != nil {
				return nil, err
			}

			proxy := httputil.NewSingleHostReverseProxy(serverUrl)
			proxy.ErrorHandler = c.getProxyErrorHandler()
			isAlive := &atomic.Bool{}
			isAlive.Store(true)

			servers = append(servers, &model.Server{
				Url:          serverUrl,
				Alive:        isAlive,
				ReverseProxy: proxy,
			})
		}
		return &model.ServerPool{
			Servers:    servers,
			CurrentIdx: 0,
		}, nil
	} else {
		return nil, errors.New("no servers attribute provided")
	}
}

func (c Config) getProxyErrorHandler() func(w http.ResponseWriter, r *http.Request, e error) {
	// TODO add retries to config file and think of other maybe
	return func(w http.ResponseWriter, r *http.Request, e error) {

	}
}
