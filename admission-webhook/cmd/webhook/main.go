package main

import (
	"flag"
	"fmt"
	"github.com/apex/log"
	"github.com/integr8ly/managed-services-controller/admission-webhook/pkg/server"
	"github.com/operator-framework/operator-sdk/pkg/k8sclient"
)

var certFile string
var keyFile string
var port string

func init() {
	flag.StringVar(&certFile, "tls-cert-file", certFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&keyFile, "tls-private-key-file", keyFile, ""+
		"File containing the default x509 private key matching --tls-cert-file.")
	flag.StringVar(&port, "port", "8443", ""+
		"Port the server will listen on --port.")
	flag.Parse()
}

func main() {
	config := server.Config{
		CertFile:  certFile,
		KeyFile:   keyFile,
		Port:      port,
		K8sClient: k8sclient.GetKubeClient(),
	}
	log.Info(fmt.Sprintf("Starting server on port: %s", config.Port))
	err := server.ListenAndServeTLS(config)
	if err != nil {
		log.Error("Error starting server: " + err.Error())
	}
}
