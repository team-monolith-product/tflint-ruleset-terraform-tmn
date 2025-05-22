# terraform_list_order

Recommends alphabetical order for variables.

## Configuration

```hcl
rule "terraform_variables_order" {
  enabled = true
  group_required = false # Set to true if you want required variables to be sorted separately from the optional ones
}
```


## Example

```hcl

variable "b" {
  type = string
}

variable "a" {
  type = string
}

```

Result:
```
$ tflint -f compact --recursive
2 issue(s) found:

main.tf:2:3: Notice - List 'names' is not sorted alphabetically. Recommended order: [Alice Bob Charlie Xavier] (terraform_list_order)
main.tf:7:5: Notice - List 'actions' is not sorted alphabetically. Recommended order: [kms:Decrypt* kms:Describe* kms:Encrypt* kms:GenerateDataKey* kms:ReEncrypt*] (terraform_list_order)
```
