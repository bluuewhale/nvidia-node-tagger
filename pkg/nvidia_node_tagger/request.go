package nvidia_node_tagger

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Request interface {
	Send() (string, error)
}

type NodePatchRequest struct {
	NodeName  string
	Clientset *kubernetes.Clientset
	Patch     *Patch
}

func (r *NodePatchRequest) Send() (*v1.Node, error) {
	logrus.Infof("Sending %s %s patch to %s", r.Patch.Op, r.Patch.Path, r.NodeName)
	data, err := json.Marshal([]Patch{*r.Patch})
	if err != nil {
		return nil, err
	}

	return r.Clientset.
		CoreV1().
		Nodes().
		Patch(context.TODO(), r.NodeName, types.JSONPatchType, data, metav1.PatchOptions{}, r.Patch.SubResources...)
}
