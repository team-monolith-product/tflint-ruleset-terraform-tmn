package rules

import (
	"fmt"
	"github.com/kenske/tflint-ruleset-terraform-sort/project"
	"reflect"
	"sort"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// TerraformVariablesOrderRule checks whether the variables are sorted in expected order
type TerraformVariablesOrderRule struct {
	tflint.DefaultRule
}

type TerraformVariablesOrderRuleConfig struct {
	SortRequired bool `hclext:"sort_required"`
}

// NewTerraformVariablesOrderRule returns a new rule
func NewTerraformVariablesOrderRule() *TerraformVariablesOrderRule {
	return &TerraformVariablesOrderRule{}
}

// Name returns the rule name
func (r *TerraformVariablesOrderRule) Name() string {
	return "terraform_variables_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVariablesOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformVariablesOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformVariablesOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are sorted in the expected order
func (r *TerraformVariablesOrderRule) Check(runner tflint.Runner) error {

	config := &TerraformVariablesOrderRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	logger.Debug(fmt.Sprintf("value of sort_required: %v", config.SortRequired))

	for _, file := range files {
		if subErr := r.checkVariablesOrder(runner, config.SortRequired, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformVariablesOrderRule) checkVariablesOrder(runner tflint.Runner, sortRequired bool, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_variables_order check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks

	var sortedVariableNames []string

	if sortRequired {
		requiredVars := r.getSortedVariableNames(blocks, true)
		optionalVars := r.getSortedVariableNames(blocks, false)
		sortedVariableNames = append(requiredVars, optionalVars...)
	} else {
		sortedVariableNames = r.getVariableNames(blocks)
		sort.Strings(sortedVariableNames)
	}

	variableNames := r.getVariableNames(blocks)
	if reflect.DeepEqual(variableNames, sortedVariableNames) {
		logger.Debug("variables are sorted")
		logger.Debug(fmt.Sprintf("variables: %v", variableNames))
		return nil
	}

	logger.Debug("variables are not sorted")
	logger.Debug(fmt.Sprintf("variables: %v", variableNames))

	firstRange := r.firstVariableRange(blocks)
	sortedVariableHclTxts := r.sortedVariableCodeTxts(blocks, file, sortedVariableNames)
	sortedVariableHclBytes := hclwrite.Format([]byte(strings.Join(sortedVariableHclTxts, "\n\n")))

	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended variables order:\n%s", sortedVariableHclBytes),
		*firstRange,
	)
}

func (r *TerraformVariablesOrderRule) sortedVariableCodeTxts(blocks hclsyntax.Blocks, file *hcl.File, sortedVariableNames []string) []string {
	variableHclTxts := r.variableCodeTxts(blocks, file)
	var sortedVariableHclTxts []string
	for _, name := range sortedVariableNames {
		sortedVariableHclTxts = append(sortedVariableHclTxts, variableHclTxts[name])
	}
	return sortedVariableHclTxts
}

func (r *TerraformVariablesOrderRule) variableCodeTxts(blocks hclsyntax.Blocks, file *hcl.File) map[string]string {
	variableHclTxts := make(map[string]string)
	r.forVariables(blocks, func(v *hclsyntax.Block) {
		name := v.Labels[0]
		variableHclTxts[name] = string(v.Range().SliceBytes(file.Bytes))
	})
	return variableHclTxts
}

func (r *TerraformVariablesOrderRule) firstVariableRange(blocks hclsyntax.Blocks) *hcl.Range {
	var firstRange *hcl.Range
	r.forVariables(blocks, func(v *hclsyntax.Block) {
		if firstRange == nil {
			firstRange = ref(v.DefRange())
		}
	})
	return firstRange
}

func (r *TerraformVariablesOrderRule) getVariableNames(blocks hclsyntax.Blocks) []string {
	var variableNames []string
	r.forVariables(blocks, func(v *hclsyntax.Block) {
		variableNames = append(variableNames, v.Labels[0])
	})
	return variableNames
}

func (r *TerraformVariablesOrderRule) getSortedVariableNames(blocks hclsyntax.Blocks, required bool) []string {
	var variableNames []string

	r.forVariables(blocks, func(v *hclsyntax.Block) {
		if _, hasDefault := v.Body.Attributes["default"]; hasDefault != required {
			variableNames = append(variableNames, v.Labels[0])
		}
	})

	sort.Strings(variableNames)
	return variableNames
}

func (r *TerraformVariablesOrderRule) forVariables(blocks hclsyntax.Blocks, action func(v *hclsyntax.Block)) {
	for _, block := range blocks {
		if block.Type == "variable" {
			action(block)
		}
	}
}
