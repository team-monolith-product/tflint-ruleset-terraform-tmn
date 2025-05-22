package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableOrderRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "No lists",
			Content: `
terraform{}`,
			Expected: helper.Issues{},
		},
		{
			Name: "Empty list",
			Content: `
locals {
  names = []
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. correct list order",
			Content: `
locals {
  names = ["Alice", "Bob", "Charlie"]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "3. Incorrect list order",
			Content: `
locals {
  names = ["Xavier", "Alice", "Bob", "Charlie"]
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformListOrderRule(),
					Message: `List 'names' is not sorted alphabetically. Recommended order: ["Alice", "Bob", "Charlie", "Xavier"]`,
				},
			},
		},
		{
			Name: "4. Correctly sorted list under nested block",
			Content: `
data "aws_iam_policy_document" "current" {
  statement {
    actions = [
      "kms:Decrypt*",
      "kms:Describe*",
      "kms:Encrypt*",
      "kms:GenerateDataKey*",
      "kms:ReEncrypt*",
    ]
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "5. Incorrectly sorted list under nested block",
			Content: `
data "aws_iam_policy_document" "current" {
  statement {
    actions = [
      "kms:Describe*",
      "kms:Encrypt*",
      "kms:GenerateDataKey*",
      "kms:ReEncrypt*",
      "kms:Decrypt*",	
    ]
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformListOrderRule(),
					Message: `List 'actions' is not sorted alphabetically. Recommended order: ["kms:Decrypt*", "kms:Describe*", "kms:Encrypt*", "kms:GenerateDataKey*", "kms:ReEncrypt*"]`,
				},
			},
		},
	}
	rule := NewTerraformListOrderRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "main.tf"
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
