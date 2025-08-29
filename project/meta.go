package project

import "fmt"

// Version is ruleset version
var Version string = "1.0.5"

// ReferenceLink returns the rule reference link
func ReferenceLink(name string) string {
	return fmt.Sprintf("https://github.com/team-monolith-product/tflint-ruleset-terraform-tmn/blob/v%s/docs/rules/%s.md", Version, name)
}
