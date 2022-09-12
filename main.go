package main

import (
	"flag"
	"github.com/QQGoblin/kube-webhook/pkg/podToleration"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	certDir string
	port    int
)

func init() {
	flag.StringVar(&certDir, "cert-dir", "",
		"CertDir is the directory that contains the server key and certificate. "+
			"The server key and certificate must be named tls.key and tls.crt, respectively.")
	flag.IntVar(&port, "port", 9443, " Port is the port that the webhook server serves at.")
	flag.Parse()
}

func main() {

	klog.Info("start up controller manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		CertDir: certDir,
		Port:    port,
	})
	if err != nil {
		klog.Fatalf("unable to set up controller manager: %+v", err)
	}

	klog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	hookServer.Register("/mutate-pod-toleration", &webhook.Admission{Handler: podToleration.NewHandler(mgr.GetClient())})

	klog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Fatalf("unable to set up controller manager: %+v", err)
	}
}
