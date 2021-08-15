package k8s

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeConfig(masterUrl, kubeConfig string) (*rest.Config, error) {
	var config *rest.Config
	var err error

	if kubeConfig == "" {
		logrus.Info("Using in-cluster kubeconfig")
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		logrus.Infof("Using out-of-cluster kubeconfig from %s", kubeConfig)

		if strings.HasPrefix(kubeConfig, "~/") {
			user, err := user.Current()
			if err != nil {
				return nil, err
			}

			kubeConfig = filepath.Join(user.HomeDir, (kubeConfig)[2:])
		}

		config, err = clientcmd.BuildConfigFromFlags(masterUrl, kubeConfig)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func NewKubeClient(config *rest.Config) (*kubernetes.Clientset, error) {
	k, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k, nil
}
