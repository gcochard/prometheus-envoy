package main

import (
	"os"
	"log"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gcochard/prometheus-proxy"
	"github.com/gcochard/prometheus-envoy/pkg"
)

var token *string
func factory(target string) prometheus.Collector {
	log.Printf("token: %s", *token)
	return pkg.NewEnvoyCollector(target, *token)
}

var app = proxy.Application {
	CreateFactory: func() proxy.CollectorFactory {
		return factory
	},
}

func main() {
	token = flag.String("token", "", "A JWT from entrez.enphaseenergy.com")
	port := flag.Int("port", 2112, "the port to listen on")
	listen := flag.String("listen", "127.0.0.1", "the address to listen on")
	flag.Parse()
	log.Printf("port: %d, listen: %s, os.args: %v", *port, *listen, os.Args)
	proxy.Main(app)
}
