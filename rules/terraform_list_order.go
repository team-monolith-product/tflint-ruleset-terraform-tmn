package rules

import (
	"fmt"
	"github.com/kenske/tflint-ruleset-terraform-sort/project"
	"reflect"
	"sort"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// TerraformListOrderRule checks whether a list is sorted in expected order
type TerraformListOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformListOrderRule returns a new rule
func NewTerraformListOrderRule() *TerraformListOrderRule {
	return &TerraformListOrderRule{}
}

// Name returns the rule name
func (r *TerraformListOrderRule) Name() string {
	return "terraform_list_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformListOrderRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformListOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformListOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are sorted in expected order
func (r *TerraformListOrderRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range body.Blocks {
			if err := r.checkBlock(runner, block, file); err != nil {
				return err
			}
		}
	}

	return nil
}
func (r *TerraformListOrderRule) checkBlock(runner tflint.Runner, block *hclsyntax.Block, file *hcl.File) error {
	// Check attributes in the current block
	for name, attr := range block.Body.Attributes {
		value, diag := attr.Expr.Value(nil)
		if diag.HasErrors() {
			logger.Debug(fmt.Sprintf("Skipping attribute '%s' due to error: %s", name, diag.Error()))
			continue
		}

		if value.Type().IsTupleType() {
			list := value.AsValueSlice()
			var items []string
			for _, v := range list {
				if v.Type().FriendlyName() != "string" {
					logger.Debug(fmt.Sprintf("Skipping attribute '%s' since it contains non-string values", name))
					continue
				}
				items = append(items, v.AsString())
			}

			if len(items) == 0 {
				continue
			}

			if isSorted(items) {
				continue
			}

			sortedItems := make([]string, len(items))
			copy(sortedItems, items)
			sort.Strings(sortedItems)
			listRange := attr.Expr.Range()
			suggestedFix := toHCLList(sortedItems)

			return runner.EmitIssueWithFix(
				r,
				fmt.Sprintf("List '%s' is not sorted alphabetically. Recommended order: %v", name, suggestedFix),
				listRange,
				func(f tflint.Fixer) error {
					err := f.ReplaceText(listRange, suggestedFix)
					if err != nil {
						return err
					}
					return nil
				},
			)

		}
	}

	// Recursively check nested blocks
	for _, nestedBlock := range block.Body.Blocks {
		if err := r.checkBlock(runner, nestedBlock, file); err != nil {
			return err
		}
	}

	return nil
}

func isSorted(list []string) bool {
	return reflect.DeepEqual(list, sortedCopy(list))
}

func sortedCopy(list []string) []string {
	sorted := make([]string, len(list))
	copy(sorted, list)
	sort.Strings(sorted)
	return sorted
}

func toHCLList(items []string) string {
	quotedItems := make([]string, len(items))
	for i, item := range items {
		quotedItems[i] = fmt.Sprintf("\"%s\"", item)
	}
	return fmt.Sprintf("[\n%s\n]", strings.Join(quotedItems, ",\n"))

}
