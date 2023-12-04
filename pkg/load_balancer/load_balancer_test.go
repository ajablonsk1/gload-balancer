package load_balancer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewLoadBalancer(t *testing.T) {
	t.Run("creates new load balancer", func(t *testing.T) {
		path := "../../config/config.json"
		lb, err := NewLoadBalancer(path)
		if err != nil {
			t.Errorf("error from new load balancer: %s", err)
		}
		if lb.addr != "localhost:8080" {
			t.Errorf("wrong load balancer address")
		}
	})
}

func TestLoadBalancer(t *testing.T) {
	t.Run("starts new load balancer", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "test")
		})
		s1 := &http.Server{
			Addr:    ":1111",
			Handler: mux,
		}
		s2 := &http.Server{
			Addr:    ":1112",
			Handler: mux,
		}
		go func() { s1.ListenAndServe() }()
		go func() { s2.ListenAndServe() }()

		lb, _ := NewLoadBalancer("../../config/config.json")
		go func() { lb.Start() }()

		time.Sleep(1 * time.Second)

		request, _ := http.NewRequest(http.MethodGet, "/test", nil)
		response := httptest.NewRecorder()
		lb.proxyHandler.ServeHTTP(response, request)

		if response.Code != 200 {
			t.Errorf("wrong response code")
		}

		if response.Body.String() != "test" {
			t.Errorf("wrong response body")
		}
	})
}
