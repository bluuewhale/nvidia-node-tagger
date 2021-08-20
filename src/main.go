package main

import (
	"context"
	"flag"
	"os"

	"github.com/BlueWhaleKo/nvidia-node-tagger/pkg/gpu"
	"github.com/BlueWhaleKo/nvidia-node-tagger/pkg/k8s"
	tagger "github.com/BlueWhaleKo/nvidia-node-tagger/pkg/nvidia_node_tagger"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ==================================
// ===== Request-line arguments =====
// ==================================
var (
	argKubecfgFile       = flag.String("kubecfg-file", "", `Location of kubecfg file for access to kubernetes master service; --kube_master_url overrides the URL part of this; if neither this nor --kube_master_url are provided, defaults to ServiceAccount tokens`)
	argKubeMasterURL     = flag.String("kube-master-url", "", `URL to reach kubernetes master. Env variables in this flag will be expanded.`)
	argNodeName          = flag.String("node", "", "name of nodes to add annotation and capacity info")
	argNamespace         = flag.String("namespace", "cluster-addons-nvidia-node-tagger", "Name of the namespace to deploy nvidia-node-taggers")
	argAnnotationsPrefix = flag.String("annotation-prefix", "nvidia-node-tagger", "prefix for node annotations")
	argLabelPrefix       = flag.String("label-prefix", "nvidia-node-tagger", "prefix for node labels")
	argCapacityPrefix    = flag.String("capacity-prefix", "nvidia-node-tagger", "prefix for node capacity")
)

func main() {
	flag.Parse()

	// parse gpu informations
	gpuInfoList, err := gpu.NewGpuDeviceList()
	if err != nil {
		logrus.Fatal(err)
	}

	// create annotation patch
	pm := tagger.NewPatchFactory(*argAnnotationsPrefix)
	patchAnnotation, err := pm.Patch("add", "/metadata/annotations", gpuInfoList)
	if err != nil {
		logrus.Fatal(err)
	}

	// create k8s client
	kubecfg, err := k8s.NewKubeConfig(*argKubeMasterURL, *argKubecfgFile)
	if err != nil {
		logrus.Fatal(err)
	}

	clientset, err := k8s.NewKubeClient(kubecfg)
	if err != nil {
		logrus.Fatal(err)
	}

	var nodeName string
	if *argNodeName != "" {
		nodeName = *argNodeName
	} else {
		// if deployed as a daemonset
		podName := os.Getenv("HOSTNAME")
		if podName == "" {
			logrus.Fatal("Environmental variable 'HOSTNAME' not found")
		}

		pod, err := clientset.CoreV1().Pods(*argNamespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			logrus.Fatal(err)
		}
		nodeName = pod.Spec.NodeName
	}

	logrus.Infof("NodeName: %s\n", nodeName)

	rq := tagger.NodePatchRequest{
		NodeName:  nodeName,
		Clientset: clientset,
		Patch:     patchAnnotation,
	}

	_, err = rq.Send()
	if err != nil {
		logrus.Fatal(err)
	}
}
