package handler

import (
	"net/http"

	"github.com/ajablonsk1/gload-balancer/internal/model"
)

type ProxyHandler struct {
	Strategy   model.LoadDistributionStrategy
	ServerPool *model.ServerPool
}

func (h ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := h.Strategy.GetServer(h.ServerPool, r.RemoteAddr)
	server.Proxy.ServeHTTP(w, r)
}
