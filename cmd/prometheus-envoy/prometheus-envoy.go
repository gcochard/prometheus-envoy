package main

import (
	"os"
	"log"
	"flag"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gcochard/prometheus-proxy"
	"github.com/gcochard/prometheus-envoy/pkg"
)

var ec *pkg.EnvoyCollector
var tokenpath *string
var token *string
func factory(target string) prometheus.Collector {
	if ec == nil{
		log.Printf("token: %s...%s", (*token)[:4], (*token)[len(*token)-4:])
		ec = pkg.NewEnvoyCollector(target, *token)
	}
	return ec
}

var app = proxy.Application {
	CreateFactory: func() proxy.CollectorFactory {
		return factory
	},
}

func main() {
	tokenpath = flag.String("token", "", "A JWT from entrez.enphaseenergy.com")
	port := flag.Int("port", 2112, "the port to listen on")
	listen := flag.String("listen", "127.0.0.1", "the address to listen on")
	flag.Parse()
	if ! strings.HasPrefix(*tokenpath, "eyJ") {
		// it's a systemd credentials file
		data, err := os.ReadFile(*tokenpath)
		if err != nil {
			log.Fatal(err)
		}
		// copy the slice into a string
		strToken := strings.TrimSpace(string(data[:]))
		// repoint token to address of strToken
		token = &strToken
	} else {
		// simple, just repoint token to tokenpath
		token = tokenpath
	}
	log.Printf("port: %d, listen: %s, os.args: %v", *port, *listen, os.Args)
	proxy.Main(app)
}
