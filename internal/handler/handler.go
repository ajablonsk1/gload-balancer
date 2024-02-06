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
	if ipHash, ok := h.Strategy.(*model.IPHash); ok {
		ipHash.CurrSourceAddress = r.RemoteAddr
	}
	server := h.Strategy.GetServer(h.ServerPool)
	server.Proxy.ServeHTTP(w, r)
}
