data "aws_iam_policy_document" "current" {

  statement {
    actions = [
      "kms:GenerateDataKey*",
      "kms:Describe*",
      "kms:Decrypt*",
      "kms:Encrypt*",
      "kms:ReEncrypt*",
      "aaa",
    ]

  }

}
