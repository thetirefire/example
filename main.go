/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/thetirefire/badidea/server"
	foov1 "github.com/thetirefire/example/api/v1"
	"github.com/thetirefire/example/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/discovery"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = foov1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func badIdeaClientConfig() clientcmd.ClientConfig {
	clustercfg := clientcmdapi.NewCluster()
	clustercfg.InsecureSkipTLSVerify = true
	clustercfg.Server = "https://localhost:6443"

	usercfg := clientcmdapi.NewAuthInfo()
	usercfg.Username = "bad"
	usercfg.Password = "idea"

	contextcfg := clientcmdapi.NewContext()
	contextcfg.AuthInfo = "badidea"
	contextcfg.Cluster = "badidea"

	apiConfig := clientcmdapi.NewConfig()
	apiConfig.Clusters["badidea"] = clustercfg
	apiConfig.AuthInfos["badidea"] = usercfg
	apiConfig.Contexts["badidea"] = contextcfg
	apiConfig.CurrentContext = "badidea"

	return clientcmd.NewDefaultClientConfig(*apiConfig, nil)
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(klogr.New())

	sigHandler := ctrl.SetupSignalHandler()
	go func() {
		if err := server.RunBadIdeaServer(sigHandler.Done()); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}()

	clientConfig := badIdeaClientConfig()

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		setupLog.Error(err, "failed to create badidea client")
		os.Exit(1)
	}

	pollContext, cancel := context.WithTimeout(sigHandler, 30*time.Second)
	defer cancel()
	if err := wait.PollUntil(
		1*time.Second,
		func() (bool, error) {
			setupLog.Info("attempting to talk to badidea apiserver")
			dc, err := discovery.NewDiscoveryClientForConfig(restConfig)
			if err != nil {
				setupLog.Error(err, "failed to create discovery client")
				return false, nil
			}
			verInfo, err := dc.ServerVersion()
			if err != nil {
				setupLog.Error(err, "failed to query server version")
				return false, nil
			}

			setupLog.Info("server responded", "version", verInfo)
			return true, nil
		},
		pollContext.Done(),
	); err != nil {
		setupLog.Error(err, "timed out attempting to create badidea client")
	}

	// TODO: deploy crd resources
	c1 := exec.Command("kustomize", "build", "config/crd")
	c2 := exec.Command("kubectl", "--server", "https://localhost:6443", "--insecure-skip-tls-verify", "--username", "bad", "--password", "idea", "create", "-f", "-")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r
	c2.Stdout = os.Stdout

	if err := c1.Start(); err != nil {
		setupLog.Error(err, "failed to start kustomize")
		os.Exit(1)
	}

	if err := c2.Start(); err != nil {
		setupLog.Error(err, "failed to start kubectl")
		os.Exit(1)
	}

	go func() {
		defer w.Close()
		if err := c1.Wait(); err != nil {
			setupLog.Error(err, "failed to run kustomize")
			os.Exit(1)
		}
	}()

	if err := c2.Wait(); err != nil {
		setupLog.Error(err, "failed to run kubectl")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "7554e382.example.thetirefire",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.BarReconciler{}).SetupWithManager(context.Background(), mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Bar")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(sigHandler); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
