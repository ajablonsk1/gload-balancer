package main

import (
	"flag"
	lb "github.com/ajablonsk1/gload-balancer/pkg/load_balancer"
)

func main() {
	configPath := flag.String("path", "", "Path to config file")
	if *configPath == "" {
		panic("You must provide flag with path to config file")
	}

	if loadBalancer, err := lb.NewLoadBalancer(*configPath); err != nil {
		panic(err.Error())
	} else {
		loadBalancer.Start()
	}
}
