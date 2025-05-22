package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariablesOrderRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "1. no variable",
			Content: `
terraform{}`,
			Config: `
rule "terraform_variables_order" {
  enabled       = true
  group_required = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. correct variable order",
			Content: `
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}`,
			Config: `
rule "terraform_variables_order" {
  enabled       = true
  group_required = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "3. sorting based on default value",
			Content: `
variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "image_id" {
  type = string
}`,
			Config: `
rule "terraform_variables_order" {
  enabled       = true
  group_required = true
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformVariablesOrderRule(),
					Message: `Recommended variables order:
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}`,
				},
			},
		},
		{
			Name: "4. sorting in alphabetic order",
			Content: `
variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}`,
			Config: `
rule "terraform_variables_order" {
  enabled       = true
  group_required = true
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformVariablesOrderRule(),
					Message: `Recommended variables order:
variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}`,
				},
			},
		},
		{
			Name: "5. mixed",
			Content: `
variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "image_id" {
  type = string
}`,
			Config: `
rule "terraform_variables_order" {
  enabled       = true
  group_required = true
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformVariablesOrderRule(),
					Message: `Recommended variables order:
variable "image_id" {
  type = string
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}`,
				},
			},
		},
		{
			Name: "6. sorting with group_required disabled",
			Content: `
variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}

variable "image_id" {
  type = string
}`,
			Config: `
rule "terraform_variables_order" {
  enabled       = true
  group_required = false
}`,
			Expected: helper.Issues{},
		},
	}
	rule := NewTerraformVariablesOrderRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "main.tf"

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content, ".tflint.hcl": tc.Config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
