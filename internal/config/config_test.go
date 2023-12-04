package config_test

import (
	"testing"

	"github.com/ajablonsk1/gload-balancer/internal/config"
	"github.com/ajablonsk1/gload-balancer/internal/model"
)

func TestGetConfig(t *testing.T) {
	t.Run("created valid config", func(t *testing.T) {
		c, err := config.GetConfig("../../config/config.json")
		if err != nil {
			t.Errorf("error while creating config: %s", err.Error())
		}

		got := c["address"]
		want := "localhost:8080"
		if got != want {
			t.Errorf("wrong address. got %s want %s", got, want)
		}

		got = c["strategy"]
		want = "round-robin"
		if got != want {
			t.Errorf("wrong stategy. got %s want %s", got, want)
		}

		got = c["servers"].([]interface{})[0].(map[string]interface{})["host"]
		want = "localhost:1111"
		if got != want {
			t.Errorf("wrong first sever host. got %s want %s", got, want)
		}
	})
}

func TestGetAdress(t *testing.T) {
	t.Run("get address from config", func(t *testing.T) {
		c, err := config.GetConfig("../../config/config.json")
		if err != nil {
			t.Errorf("error while creating config: %s", err.Error())
		}

		address, err := c.GetAddress()
		if err != nil {
			t.Errorf("error getting address: %s", err.Error())
		}

		got := address
		want := "localhost:8080"
		if got != want {
			t.Errorf("wrong address. got %s want %s", got, want)
		}
	})
}

func TestGetStrategy(t *testing.T) {
	t.Run("get strategy from config", func(t *testing.T) {
		c, err := config.GetConfig("../../config/config.json")
		if err != nil {
			t.Errorf("error while creating config: %s", err.Error())
		}

		strategy, err := c.GetLoadStrategy()
		if err != nil {
			t.Errorf("error getting load strategy: %s", err.Error())
		}

		got := strategy
		want := &model.RoundRobin{}
		if _, ok := strategy.(*model.RoundRobin); !ok {
			t.Errorf("wrong strategy. got %T want %T", got, want)
		}
	})
}

func TestGetServerPool(t *testing.T) {
	t.Run("get server pool from config", func(t *testing.T) {
		c, err := config.GetConfig("../../config/config.json")
		if err != nil {
			t.Errorf("error while creating config: %s", err.Error())
		}

		serverPool, err := c.GetServerPool()
		if err != nil {
			t.Errorf("error getting server pool: %s", err.Error())
		}

		got := serverPool.Servers[0].Url.String()
		want := "http://localhost:1111"
		if got != want {
			t.Errorf("wrong first server in server pool. got %s want %s", got, want)
		}
	})
}
