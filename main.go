package main

import (
	"github.com/kenske/tflint-ruleset-terraform-sort/project"
	"github.com/kenske/tflint-ruleset-terraform-sort/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "terraform-sort",
			Version: project.Version,
			Rules: []tflint.Rule{
				rules.NewTerraformListOrderRule(),
				rules.NewTerraformVariablesOrderRule(),
				rules.NewTerraformResourceOrderRule(),
			},
		},
	})
}
