package nvidia_node_tagger

import (
	"fmt"

	"github.com/fatih/structs"
)

func Map(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

// Flatten transform nested map into a flat map
func Flatten(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}

// FlattenMap transform an interface into a flat map
func FlattenMap(s interface{}) map[string]interface{} {
	m := Map(s)
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		default:
			o[k] = v
		}
	}
	return o
}

// AddPrefix to keys in a map
func AddPrefix(m *map[string]interface{}, prefix string) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range *m {
		output[fmt.Sprintf("%s/%s", prefix, k)] = fmt.Sprintf("%v", v)
	}
	return output
}
