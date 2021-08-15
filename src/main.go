package main

import (
	"context"
	"flag"
	"os"

	"github.com/BlueWhaleKo/nvidia-node-tagger/pkg/gpu"
	"github.com/BlueWhaleKo/nvidia-node-tagger/pkg/k8s"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ==================================
// ===== Command-line arguments =====
// ==================================
var (
	flags            = flag.NewFlagSet("", flag.ContinueOnError)
	argKubecfgFile   = flags.String("kubecfg-file", "", `Location of kubecfg file for access to kubernetes master service; --kube_master_url overrides the URL part of this; if neither this nor --kube_master_url are provided, defaults to ServiceAccount tokens`)
	argKubeMasterURL = flags.String("kube-master-url", "", `URL to reach kubernetes master. Env variables in this flag will be expanded.`)
	argNamespace     = flags.String("namespace", "cluster-addons-gpu-node-tagger", "Name of the namespace to deploy gpu-node-taggers")
	argLabelPrefix   = flags.String("label-prefix", "BlueWhaleKo.com", "prefix for node labels")
)

func parseArgs() {
	flags.Parse(os.Args)
}

func main() {
	parseArgs()

	gpuStatList, err := gpu.NewGpuStatList()
	if err != nil {
		logrus.Fatal(err)
	}

	kubecfg, err := k8s.NewKubeConfig(*argKubeMasterURL, *argKubecfgFile)
	if err != nil {
		logrus.Fatal(err)
	}

	clientset, err := k8s.NewKubeClient(kubecfg)
	if err != nil {
		logrus.Fatal(err)
	}

	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		logrus.Fatal("Environmental variable 'HOSTNAME' not found")
	}

	pod, err := clientset.CoreV1().Pods("").Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		logrus.Fatal(err)
	}
	nodeName := pod.Spec.NodeName

	patch := k8s.NewPatchAddLabels(map[string]interface{}{
		"gpu.memory.total.sum": gpuStatList.MemoryTotalSum,
		"gpu.memory.used.sum":  gpuStatList.MemoryUsedSum,
		"gpu.memory.free.sum":  gpuStatList.MemoryFreeSum,
	})

	command := k8s.NodePatchCommand{
		NodeName:  nodeName,
		Clientset: clientset,
		Patch:     patch,
	}

	command.Execute()
}
