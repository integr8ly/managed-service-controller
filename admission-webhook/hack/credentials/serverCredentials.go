//TODO: What is the story with these package names?
package main

import (
	"crypto/x509"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/util/cert"
	"net"
	"os"
	"path/filepath"
)

type certContext struct {
	serverCert []byte
	serverkey  []byte
	caCert     []byte
}

func createCredentials(config cert.Config) *certContext {
	caKey, err := cert.NewPrivateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create CA private serverKey %v", err)
	}

	caCert, err := cert.NewSelfSignedCACert(cert.Config{CommonName: "server-cert-ca"}, caKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create CA cert for apiserver %v", err)
	}

	serverKey, err := cert.NewPrivateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create private serverKey for %v", err)
	}

	serverCert, err := cert.NewSignedCert(config, serverKey, caCert, caKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create cert%v", err)
	}

	return &certContext{
		serverCert: cert.EncodeCertPEM(serverCert),
		serverkey:  cert.EncodePrivateKeyPEM(serverKey),
		caCert:     cert.EncodeCertPEM(caCert),
	}
}

func createLocalCredentials(ipAddress string) *certContext {
	return createCredentials(cert.Config{
		CommonName: ipAddress,
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		AltNames: cert.AltNames{
			DNSNames: []string{
				"localhost",
			},
			IPs: []net.IP{
				net.ParseIP(ipAddress),
				net.ParseIP("127.0.0.1"),
				net.ParseIP("::1"),
			},
		},
	})
}

func generateLocalCredentials(ipAddress string) {
	context := createLocalCredentials(ipAddress)

	fmt.Fprintln(os.Stdout, "Base64 encoded CA cert")
	fmt.Fprintln(os.Stdout, b64.URLEncoding.EncodeToString(context.caCert))
	fmt.Fprintln(os.Stdout, "\n")

	keyPath := filepath.Join(".", "build", "tmp", "server-credentials")
	os.MkdirAll(keyPath, os.ModePerm)

	if err := ioutil.WriteFile("./build/tmp/server-credentials/server.key.pem", context.serverkey, 0644); err != nil {
		fmt.Fprintln(os.Stdout, "Error writing server key to file")
		return
	}
	fmt.Fprintln(os.Stdout, "Server key is at ./build/tmp/server-credentials/server.key.pem")

	if err := ioutil.WriteFile("./build/tmp/server-credentials/server.cert.pem", context.serverCert, 0600); err != nil {
		fmt.Fprintln(os.Stdout, "Error writing server key to file")
		return
	}
	fmt.Fprintln(os.Stdout, "Server cert is at ./build/tmp/server-credentials/server.cert.pem")
}

func createServiceCredentials(serviceName, namespace string) *certContext {
	return createCredentials(cert.Config{
		CommonName: serviceName + "." + namespace + ".svc",
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	})
}

func generateServiceCredentials(serviceName, namespace string) {
	context := createServiceCredentials(serviceName, namespace)

	fmt.Fprintln(os.Stdout, "Base64 encoded CA cert")
	fmt.Fprintln(os.Stdout, b64.URLEncoding.EncodeToString(context.caCert))
	fmt.Fprintln(os.Stdout, "\n")
	fmt.Fprintln(os.Stdout, "Base64 encoded Server key")
	fmt.Fprintln(os.Stdout, b64.URLEncoding.EncodeToString(context.serverkey))
	fmt.Fprintln(os.Stdout, "\n")
	fmt.Fprintln(os.Stdout, "Base64 encoded Server cert")
	fmt.Fprintln(os.Stdout, b64.URLEncoding.EncodeToString(context.serverCert))
}
