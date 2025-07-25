# tflint-ruleset-terraform-sort

## Requirements

- TFLint v0.42+
- Go v1.24

## Installation

You can install the plugin with `tflint --init`. Declare a config in `.tflint.hcl` as follows:

```hcl
plugin "terraform-sort" {
  enabled = true

  version = "0.3.0"
  source  = "github.com/kenske/tflint-ruleset-terraform-sort"

}
```

## Rules

See [docs](docs) for a list of rules available in this ruleset.

## Development

Clone the repository locally and run the following command:

```
$ make install
```

This will install the plugin in `~/.tflint.d/plugins`.

Run `TFLINT_LOG=debug tflint` to see debug logs.