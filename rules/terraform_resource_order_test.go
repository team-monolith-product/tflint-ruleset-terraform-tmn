package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformResourceOrderRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "1. no resources",
			Content: `
terraform{}`,
			Config: `
rule "terraform_resource_order" {
  enabled = true
  group_by_type = false
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "2. single resource",
			Content: `
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
			Config: `
rule "terraform_resource_order" {
  enabled = true
  group_by_type = false
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "3. resources in alphabetical order",
			Content: `
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_instance" "database" {
  ami           = "ami-87654321"
  instance_type = "t3.medium"
}

resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}`,
			Config: `
rule "terraform_resource_order" {
  enabled = true
  group_by_type = false
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "4. resources not in alphabetical order",
			Content: `
resource "aws_instance" "database" {
  ami           = "ami-87654321"
  instance_type = "t3.medium"
}

resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
			Config: `
rule "terraform_resource_order" {
  enabled = true
  group_by_type = false
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceOrderRule(),
					Message: `Recommended resource order:
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_instance" "database" {
  ami           = "ami-87654321"
  instance_type = "t3.medium"
}`,
				},
			},
		},
		{
			Name: "5. mixed resources not in order",
			Content: `
resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}

resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
			Config: `
rule "terraform_resource_order" {
  enabled = true
  group_by_type = false
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceOrderRule(),
					Message: `Recommended resource order:
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}`,
				},
			},
		},
		{
			Name: "6. resources grouped by type - correct order",
			Content: `
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_instance" "database" {
  ami           = "ami-87654321"
  instance_type = "t3.medium"
}

resource "aws_s3_bucket" "backup" {
  bucket = "my-backup-bucket"
}

resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}`,
			Config: `
rule "terraform_resource_order" {
  enabled = true
  group_by_type = true
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "7. resources grouped by type - incorrect order",
			Content: `
resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}

resource "aws_instance" "database" {
  ami           = "ami-87654321"
  instance_type = "t3.medium"
}

resource "aws_s3_bucket" "backup" {
  bucket = "my-backup-bucket"
}

resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
			Config: `
rule "terraform_resource_order" {
  enabled = true
  group_by_type = true
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceOrderRule(),
					Message: `Recommended resource order:
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_instance" "database" {
  ami           = "ami-87654321"
  instance_type = "t3.medium"
}

resource "aws_s3_bucket" "backup" {
  bucket = "my-backup-bucket"
}

resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}`,
				},
			},
		},
	}

	rule := NewTerraformResourceOrderRule()

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