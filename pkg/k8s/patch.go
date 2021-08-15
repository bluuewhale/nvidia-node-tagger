package k8s

import (
	"context"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Patch struct {
	Op    string                 `json:"op"`
	Path  string                 `json:"path"`
	Value map[string]interface{} `json:"value"`
}

func NewPatch(op, path string, value map[string]interface{}) Patch {
	return Patch{
		Op:    op,
		Path:  path,
		Value: value,
	}
}

func NewPatchAddLabels(value map[string]interface{}) Patch {
	return Patch{
		Op:    "add",
		Path:  "/metadata/label/",
		Value: value,
	}
}
func NewPatchReplaceLabels(value map[string]interface{}) Patch {
	return Patch{
		Op:    "replace",
		Path:  "/metadata/labels/",
		Value: value,
	}
}

type Command interface {
	Execute() (string, error)
}

type NodePatchCommand struct {
	NodeName  string
	Clientset *kubernetes.Clientset
	Patch     Patch
}

func (p *NodePatchCommand) Execute() (*v1.Node, error) {
	data, err := json.Marshal(p.Patch)
	if err != nil {
		return nil, err
	}

	return p.Clientset.CoreV1().Nodes().Patch(context.TODO(), p.NodeName, types.JSONPatchType, data, metav1.PatchOptions{})
}
