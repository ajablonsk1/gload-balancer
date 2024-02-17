package model_test

import (
	"testing"

	"github.com/ajablonsk1/gload-balancer/internal/config"
	"github.com/ajablonsk1/gload-balancer/internal/model"
	"github.com/ajablonsk1/gload-balancer/internal/utils"
)

func TestRoundRobinGetServer(t *testing.T) {
	t.Run("get proper server from server pool round robin", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.RoundRobin{}

		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		next := strategy.GetServer(serverPool, "localhost:22224")

		got := next.Url.String()
		want := "http://localhost:1112"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestRoundRobinGetServerStickySession(t *testing.T) {
	t.Run("get proper server from server pool round robin", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.RoundRobin{}

		s := strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22222")
		next := strategy.GetServer(serverPool, "localhost:22222")

		got := next.Url.String()
		want := s.Url.String()
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

		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		next := strategy.GetServer(serverPool, "localhost:22224")

		got := next.Url.String()
		want := "http://localhost:1111"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}

		next = strategy.GetServer(serverPool, "localhost:22225")
		got = next.Url.String()
		want = "http://localhost:1112"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestWeightedRoundRobinGetServerStickySession(t *testing.T) {
	t.Run("get proper server from server pool weighted round robin", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.WeightedRoundRobin{}

		s := strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		next := strategy.GetServer(serverPool, "localhost:22222")

		got := next.Url.String()
		want := s.Url.String()
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}

		next = strategy.GetServer(serverPool, "localhost:22225")
		got = next.Url.String()
		want = "http://localhost:1111"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}

		next = strategy.GetServer(serverPool, "localhost:22226")
		got = next.Url.String()
		want = "http://localhost:1112"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestIpHashGetServer(t *testing.T) {
	t.Run("get proper server from server pool weighted round robin", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		addr := "localhost:2121"
		strategy := model.IPHash{}
		hash := utils.Hash(addr)

		got := strategy.GetServer(serverPool, addr).Url.String()
		want := serverPool.Servers[hash%uint32(len(serverPool.Servers))].Url.String()
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestIpHashGetServerStickySession(t *testing.T) {
	t.Run("get proper server from server pool weighted round robin", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		addr := "localhost:2121"
		strategy := model.IPHash{}

		s := strategy.GetServer(serverPool, addr)
		got := strategy.GetServer(serverPool, addr).Url.String()
		want := s.Url.String()
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestLeastSessionsGetServer(t *testing.T) {
	t.Run("get proper server from server pool least connections", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.LeastSession{}

		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		_ = strategy.GetServer(serverPool, "localhost:22224")
		next := strategy.GetServer(serverPool, "localhost:22225")

		got := next.Url.String()
		want := "http://localhost:1111"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestLeastSessionsGetServerStickySession(t *testing.T) {
	t.Run("get proper server from server pool least connections", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.LeastSession{}

		s := strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		_ = strategy.GetServer(serverPool, "localhost:22222")
		next := strategy.GetServer(serverPool, "localhost:22222")

		got := next.Url.String()
		want := s.Url.String()
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestLeastSessionsGetServerStickySessionExtended(t *testing.T) {
	t.Run("get proper server from server pool least connections", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.LeastSession{}

		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22224")
		next := strategy.GetServer(serverPool, "localhost:22225")

		got := next.Url.String()
		want := "http://localhost:1111"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestLeastSessionsGetServerStickySessionFailover(t *testing.T) {
	t.Run("get proper server from server pool least connections, failover test", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.LeastSession{}

		s := strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22224")
		s.SetAlive(false)
		next := strategy.GetServer(serverPool, "localhost:22225")

		got := next.Url.String()
		want := "http://localhost:1112"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
		if s.NumberOfStickySessions() != 0 && len(s.StickySessions) != 0 {
			t.Errorf("sticky session not cleared")
		}
	})
}

func TestWeightedLeastSessionsGetServer(t *testing.T) {
	t.Run("get proper server from server pool least connections", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.WeightedLeastSession{}

		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		next := strategy.GetServer(serverPool, "localhost:22225")

		got := next.Url.String()
		want := "http://localhost:1111"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestWeightedLeastSessionsGetServerExtended(t *testing.T) {
	t.Run("get proper server from server pool least connections", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.WeightedLeastSession{}

		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		_ = strategy.GetServer(serverPool, "localhost:22224")
		_ = strategy.GetServer(serverPool, "localhost:22225")
		_ = strategy.GetServer(serverPool, "localhost:22226")
		next := strategy.GetServer(serverPool, "localhost:22227")

		got := next.Url.String()
		want := "http://localhost:1112"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}

func TestWeightedLeastSessionsGetServerStickySession(t *testing.T) {
	t.Run("get proper server from server pool least connections", func(t *testing.T) {
		c, _ := config.GetConfig("../../config/config.json")
		serverPool, _ := c.GetServerPool()
		strategy := model.WeightedLeastSession{}

		_ = strategy.GetServer(serverPool, "localhost:22222")
		_ = strategy.GetServer(serverPool, "localhost:22223")
		_ = strategy.GetServer(serverPool, "localhost:22224")
		_ = strategy.GetServer(serverPool, "localhost:22225")
		_ = strategy.GetServer(serverPool, "localhost:22222")
		next := strategy.GetServer(serverPool, "localhost:22227")

		got := next.Url.String()
		want := "http://localhost:1111"
		if got != want {
			t.Errorf("wrong server. got %s want %s", got, want)
		}
	})
}
