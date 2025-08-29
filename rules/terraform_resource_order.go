package rules

import (
	"fmt"
	"github.com/team-monolith-product/tflint-ruleset-terraform-tmn/project"
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

type TerraformResourceOrderRule struct {
	tflint.DefaultRule
}

type TerraformResourceOrderRuleConfig struct {
	GroupByType bool `hclext:"group_by_type,optional"`
}

func NewTerraformResourceOrderRule() *TerraformResourceOrderRule {
	return &TerraformResourceOrderRule{}
}

func (r *TerraformResourceOrderRule) Name() string {
	return "terraform_resource_order"
}

func (r *TerraformResourceOrderRule) Enabled() bool {
	return false
}

func (r *TerraformResourceOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformResourceOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *TerraformResourceOrderRule) Check(runner tflint.Runner) error {
	config := &TerraformResourceOrderRuleConfig{
		GroupByType: true,
	}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		if subErr := r.checkResourceOrder(runner, config.GroupByType, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformResourceOrderRule) checkResourceOrder(runner tflint.Runner, groupByType bool, file *hcl.File) error {
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		logger.Debug("skip terraform_resource_order check since it's not hcl file")
		return nil
	}
	blocks := body.Blocks

	resources := r.getResourceBlocks(blocks)
	if len(resources) == 0 {
		return nil
	}

	var sortedResourceKeys []string
	if groupByType {
		sortedResourceKeys = r.getSortedResourceKeysGrouped(resources)
	} else {
		sortedResourceKeys = r.getSortedResourceKeys(resources)
	}

	currentResourceKeys := r.getCurrentResourceKeys(resources)
	if reflect.DeepEqual(currentResourceKeys, sortedResourceKeys) {
		return nil
	}

	firstRange := r.firstResourceRange(resources)
	sortedResourceHclTxts := r.sortedResourceCodeTxts(resources, file, sortedResourceKeys)
	sortedResourceHclBytes := hclwrite.Format([]byte(strings.Join(sortedResourceHclTxts, "\n\n")))

	return runner.EmitIssueWithFix(
		r,
		fmt.Sprintf("Recommended resource order:\n%s", sortedResourceHclBytes),
		*firstRange,
		func(f tflint.Fixer) error {
			return r.applyFix(f, resources, sortedResourceKeys, file)
		},
	)
}

func (r *TerraformResourceOrderRule) getResourceBlocks(blocks hclsyntax.Blocks) []*hclsyntax.Block {
	var resources []*hclsyntax.Block
	for _, block := range blocks {
		if block.Type == "resource" && len(block.Labels) >= 2 {
			resources = append(resources, block)
		}
	}
	return resources
}

func (r *TerraformResourceOrderRule) getResourceKey(block *hclsyntax.Block) string {
	if len(block.Labels) >= 2 {
		return fmt.Sprintf("%s.%s", block.Labels[0], block.Labels[1])
	}
	return ""
}

func (r *TerraformResourceOrderRule) getCurrentResourceKeys(resources []*hclsyntax.Block) []string {
	var keys []string
	for _, resource := range resources {
		keys = append(keys, r.getResourceKey(resource))
	}
	return keys
}

func (r *TerraformResourceOrderRule) getSortedResourceKeys(resources []*hclsyntax.Block) []string {
	keys := r.getCurrentResourceKeys(resources)
	sort.Strings(keys)
	return keys
}

func (r *TerraformResourceOrderRule) getSortedResourceKeysGrouped(resources []*hclsyntax.Block) []string {
	typeMap := make(map[string][]string)
	var types []string

	for _, resource := range resources {
		resourceType := resource.Labels[0]
		key := r.getResourceKey(resource)

		if _, exists := typeMap[resourceType]; !exists {
			types = append(types, resourceType)
			typeMap[resourceType] = []string{}
		}
		typeMap[resourceType] = append(typeMap[resourceType], key)
	}

	sort.Strings(types)

	var sortedKeys []string
	for _, resourceType := range types {
		sort.Strings(typeMap[resourceType])
		sortedKeys = append(sortedKeys, typeMap[resourceType]...)
	}

	return sortedKeys
}

func (r *TerraformResourceOrderRule) firstResourceRange(resources []*hclsyntax.Block) *hcl.Range {
	if len(resources) > 0 {
		return ref(resources[0].DefRange())
	}
	return nil
}

func (r *TerraformResourceOrderRule) sortedResourceCodeTxts(resources []*hclsyntax.Block, file *hcl.File, sortedKeys []string) []string {
	resourceHclTxts := r.resourceCodeTxts(resources, file)
	var sortedResourceHclTxts []string
	for _, key := range sortedKeys {
		if txt, exists := resourceHclTxts[key]; exists {
			sortedResourceHclTxts = append(sortedResourceHclTxts, txt)
		}
	}
	return sortedResourceHclTxts
}

func (r *TerraformResourceOrderRule) resourceCodeTxts(resources []*hclsyntax.Block, file *hcl.File) map[string]string {
	resourceHclTxts := make(map[string]string)
	for _, resource := range resources {
		key := r.getResourceKey(resource)
		resourceHclTxts[key] = string(resource.Range().SliceBytes(file.Bytes))
	}
	return resourceHclTxts
}

func (r *TerraformResourceOrderRule) applyFix(f tflint.Fixer, resources []*hclsyntax.Block, sortedResourceKeys []string, file *hcl.File) error {
	if len(resources) == 0 {
		return nil
	}

	// Create a map of resource keys to blocks for easy lookup
	resourceMap := make(map[string]*hclsyntax.Block)
	for _, resource := range resources {
		key := r.getResourceKey(resource)
		resourceMap[key] = resource
	}

	// Get the text content for each resource
	resourceTexts := r.resourceCodeTxts(resources, file)

	// Build the sorted text
	var sortedTexts []string
	for _, key := range sortedResourceKeys {
		if text, exists := resourceTexts[key]; exists {
			sortedTexts = append(sortedTexts, text)
		}
	}

	// Calculate the range from the first to the last resource
	firstResource := resources[0]
	lastResource := resources[len(resources)-1]
	
	replaceRange := hcl.Range{
		Filename: firstResource.Range().Filename,
		Start:    firstResource.Range().Start,
		End:      lastResource.Range().End,
	}

	// Replace all resources with the sorted version
	sortedContent := strings.Join(sortedTexts, "\n\n")
	return f.ReplaceText(replaceRange, sortedContent)
}
