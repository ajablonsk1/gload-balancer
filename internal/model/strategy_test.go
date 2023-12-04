package model_test

import (
	"testing"

	"github.com/ajablonsk1/gload-balancer/internal/config"
	"github.com/ajablonsk1/gload-balancer/internal/model"
)

func TestRoundRobinGetServer(t *testing.T) {
	t.Run("get proper server from server pool round robin", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.RoundRobin{}

		_ = strategy.GetServer(serverPool)
		_ = strategy.GetServer(serverPool)
		next := strategy.GetServer(serverPool)

		got := next.Url.String()
		want := "http://localhost:1112"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestWeightedRoundRobinGetServer(t *testing.T) {
	t.Run("get proper server from server pool weighted round robin", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.WeightedRoundRobin{}

		_ = strategy.GetServer(serverPool)
		_ = strategy.GetServer(serverPool)
		next := strategy.GetServer(serverPool)

		got := next.Url.String()
		want := "http://localhost:1111"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}

		next = strategy.GetServer(serverPool)
		got = next.Url.String()
		want = "http://localhost:1112"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}
