package server

import (
	"crypto/tls"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
)

type Config struct {
	CertFile  string
	KeyFile   string
	Port      string
	K8sClient kubernetes.Interface
}

func configTLS(config Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		glog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
}
