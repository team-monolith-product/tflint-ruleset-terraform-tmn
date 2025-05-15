# terraform_list_order

Recommends proper order for lists. Should work within nested blocks.

## Example

```hcl

locals {
  names = ["Xavier", "Alice", "Bob", "Charlie"]
}

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
}
```

Result:
```
$ tflint -f compact --recursive
2 issue(s) found:

main.tf:2:3: Notice - List 'names' is not sorted alphabetically. Recommended order: [Alice Bob Charlie Xavier] (terraform_list_order)
main.tf:7:5: Notice - List 'actions' is not sorted alphabetically. Recommended order: [kms:Decrypt* kms:Describe* kms:Encrypt* kms:GenerateDataKey* kms:ReEncrypt*] (terraform_list_order)
```
