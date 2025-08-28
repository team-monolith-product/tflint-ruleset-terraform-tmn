# terraform_resource_order

Ensures that Terraform resources are sorted in alphabetical order.

## Configuration

```hcl
rule "terraform_resource_order" {
  enabled = true
  group_by_type = false  # Optional: group resources by type first, then sort alphabetically within each type
}
```

## Options

- `group_by_type` (bool): When set to `true`, resources are first grouped by their type (e.g., all `aws_instance` resources together), then sorted alphabetically within each group. Default is `true`.

## Examples

### Incorrect (default alphabetical sorting)

```hcl
resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}

resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}
```

### Correct (default alphabetical sorting)

```hcl
resource "aws_instance" "app" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}

resource "aws_s3_bucket" "storage" {
  bucket = "my-storage-bucket"
}
```

### Correct (with group_by_type enabled)

```hcl
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
}
```

## Why

Maintaining a consistent ordering of resources improves code readability and makes it easier to locate specific resources in large Terraform files. Alphabetical ordering provides a predictable structure that helps teams collaborate more effectively.

## How to Fix

Reorder your resources alphabetically by their full identifier (type.name). If `group_by_type` is enabled, first group resources by type, then sort each group alphabetically by resource name.