package project

import "fmt"

// Version is ruleset version
var Version string = "0.2.0"

// ReferenceLink returns the rule reference link
func ReferenceLink(name string) string {
	return fmt.Sprintf("https://github.com/kenske/tflint-ruleset-terraform-sort/blob/v%s/docs/rules/%s.md", Version, name)
}
