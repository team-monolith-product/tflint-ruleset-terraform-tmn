package main

import (
	"github.com/team-monolith-product/tflint-ruleset-terraform-tmn/project"
	"github.com/team-monolith-product/tflint-ruleset-terraform-tmn/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "terraform-tmn",
			Version: project.Version,
			Rules: []tflint.Rule{
				rules.NewTerraformListOrderRule(),
				rules.NewTerraformVariablesOrderRule(),
				rules.NewTerraformResourceOrderRule(),
			},
		},
	})
}
