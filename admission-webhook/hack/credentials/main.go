package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var local bool
var ipAddress string

var service bool
var serviceName string
var namespace string

func init() {
	flag.BoolVar(&local, "local", false, "Create certs for a locally running server --local")
	flag.StringVar(&ipAddress, "ip-address", "", "The IP Adress of your locally running server --ip-address")

	flag.BoolVar(&service, "service", false, "Create certs for a server running in a cluster with a service -service")
	flag.StringVar(&serviceName, "service-name", "", "The name of the service your server is fronted by --service-name")
	flag.StringVar(&namespace, "namespace", "", "The name of the namespace the server will be deployed to --namespace")

	flag.Parse()
}

func main() {
	if local {
		ip := net.ParseIP(ipAddress)
		if ip == nil {
			fmt.Fprintln(os.Stdout, "--ip-address should be a valid IP Address")
			return
		}

		generateLocalCredentials(ipAddress)
	}

	if service {
		if len(serviceName) == 0 || len(namespace) == 0 {
			fmt.Fprintln(os.Stdout, "--service-name and --namespace must be defined")
			return
		}

		generateServiceCredentials(serviceName, namespace)
	}
}
