package nvidia_node_tagger

import (
	"context"
	"encoding/json"
	"fmt"

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

func (p *NodePatchRequest) Send() (*v1.Node, error) {

	for k, v := range p.Patch.Value {
		v = fmt.Sprintf("%v", v)
		p.Patch.Value[k] = v
		logrus.Infof("%s: %v", k, v)
	}

	data, err := json.Marshal([]Patch{*p.Patch})
	if err != nil {
		return nil, err
	}

	return p.Clientset.
		CoreV1().
		Nodes().
		Patch(context.TODO(), p.NodeName, types.JSONPatchType, data, metav1.PatchOptions{})
}
