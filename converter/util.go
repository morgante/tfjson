package converter

import (
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

type output map[string]interface{}

func insert(out output, path []string, key string, value interface{}) {
	if len(path) > 0 && path[0] == "root" {
		path = path[1:]
	}
	for _, elem := range path {
		switch nested := out[elem].(type) {
		case output:
			out = nested
		default:
			new := output{}
			out[elem] = new
			out = new
		}
	}
	out[key] = value
}

func convertModuleDiff(out output, diff *terraform.ModuleDiff, flat bool) {
	var flatName strings.Builder
	for k, v := range diff.Resources {
		var instanceName []string
		if flat {
			flatName.Reset()

			if len(diff.Path) > 1 {
				// slice [1:] to omit "Root" string
				flatName.WriteString(strings.Join(diff.Path[1:], "."))
				flatName.WriteString(".")
			}

			flatName.WriteString(k)
			instanceName = []string{flatName.String()}
		} else {
			insert(out, diff.Path, "destroy", diff.Destroy)
			instanceName = append(diff.Path, k)
		}
		convertInstanceDiff(out, instanceName, v)
	}
}

func convertInstanceDiff(out output, path []string, diff *terraform.InstanceDiff) {
	insert(out, path, "destroy", diff.Destroy)
	insert(out, path, "destroy_tainted", diff.DestroyTainted)
	for k, v := range diff.Attributes {
		insert(out, path, k, v.New)
	}
}
