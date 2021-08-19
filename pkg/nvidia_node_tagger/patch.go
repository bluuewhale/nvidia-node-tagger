package nvidia_node_tagger

import "fmt"

type Patch struct {
	Op    string                 `json:"op"`
	Path  string                 `json:"path"`
	Value map[string]interface{} `json:"value"`
}

type PatchFactory struct {
	prefix string
}

func NewPatchFactory(prefix string) *PatchFactory {
	return &PatchFactory{
		prefix: prefix,
	}
}

func (pf *PatchFactory) Patch(op, path string, value interface{}) *Patch {
	json := FlattenMap(value)

	o := make(map[string]interface{})
	for k, v := range json {
		if pf.prefix != "" {
			key := fmt.Sprintf("%s/%s", pf.prefix, k)
			o[key] = v
		} else {
			o[k] = v
		}
	}

	return &Patch{
		Op:    op,
		Path:  path,
		Value: o,
	}
}
