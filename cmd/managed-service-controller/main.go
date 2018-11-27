package main

import (
	"context"
	"github.com/gobuffalo/packr"
	"github.com/integr8ly/managed-service-controller/pkg/handlers"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"runtime"
	"time"
	"encoding/json"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()

	sdk.ExposeMetricsPort()

	config := packr.NewBox("../../config")
	sCfgBytes, err := config.Find("service-config.json"); if err != nil {
		logrus.Fatalf("failed to get managed service config: %v", err)
	}

	var sCfg map[string]map[string]string
	json.Unmarshal(sCfgBytes, &sCfg)

	resource := "integreatly.org/v1alpha1"
	kind := "ManagedServiceNamespace"
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		logrus.Fatalf("failed to get watch namespace: %v", err)
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	k8sCfg, err := kubeConfig.ClientConfig()
	if err != nil {
		logrus.Fatalf("Error creating kube client config: %v", err)
	}

	resyncPeriod := time.Duration(5) * time.Second
	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)
	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(handlers.NewHandler(k8sCfg, sCfg))
	sdk.Run(context.TODO())
}
