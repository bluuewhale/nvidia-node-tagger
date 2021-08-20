package nvidia_node_tagger

import "fmt"

type Patch struct {
	Op    string            `json:"op"`
	Path  string            `json:"path"`
	Value map[string]string `json:"value"`
}

type PatchBuilder struct {
	Op     string
	Path   string
	Value  map[string]interface{}
	Prefix string
}

func NewPatchBuilder() *PatchBuilder {
	return &PatchBuilder{}
}

func (b *PatchBuilder) WithOperation(op string) *PatchBuilder {
	b.Op = op

	return b
}

func (b *PatchBuilder) WithPath(path string) *PatchBuilder {
	b.Path = path

	return b
}

func (b *PatchBuilder) WithValue(value map[string]interface{}) *PatchBuilder {
	b.Value = value

	return b
}

func (b *PatchBuilder) WithPrefix(prefix string) *PatchBuilder {
	b.Prefix = prefix

	return b
}

func (b *PatchBuilder) Inspect() error {
	if b.Op == "" {
		return fmt.Errorf("Operation must be set with WithOperation()")
	}

	if b.Path == "" {
		return fmt.Errorf("Path must be set with WithPath()")
	}

	if b.Value == nil {
		return fmt.Errorf("Value must be set with WithValue()")
	}

	return nil
}

func (b *PatchBuilder) Build() (*Patch, error) {

	if err := b.Inspect(); err != nil {
		return nil, err
	}

	o, err := FlattenMap(b.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to create Patch. %s", err)
	}

	value := make(map[string]string)
	for k, v := range o {

		if b.Prefix != "" {
			k = fmt.Sprintf("%s/%s", b.Prefix, k)
		}
		v := fmt.Sprintf("%v", v)

		value[k] = v
	}

	return &Patch{
		Op:    b.Op,
		Path:  b.Path,
		Value: value,
	}, nil
}
