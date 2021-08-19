package nvidia_node_tagger

import (
	"context"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Command interface {
	Execute() (string, error)
}

type NodePatchCommand struct {
	NodeName  string
	Clientset *kubernetes.Clientset
	Patch     *Patch
}

func (p *NodePatchCommand) Execute() (*v1.Node, error) {
	data, err := json.Marshal([]Patch{*p.Patch})
	if err != nil {
		return nil, err
	}

	return p.Clientset.CoreV1().Nodes().Patch(context.TODO(), p.NodeName, types.JSONPatchType, data, metav1.PatchOptions{})
}
