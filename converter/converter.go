package converter

import (
	"os"

	"github.com/hashicorp/terraform/terraform"
)

func ConvertPlan(planfile string, flatten bool) (output, error) {
	f, err := os.Open(planfile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		return nil, err
	}

	diff := output{}
	for _, v := range plan.Diff.Modules {
		convertModuleDiff(diff, v, flatten)
	}

	return diff, nil
}
