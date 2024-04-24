---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "functions Provider"
description: |-
  The functions provider extends the available functions for a Terraform project.
---

# functions Provider
The functions provider extends the available functions for a Terraform project.

## Example Usage

**versions.tf**
```terraform
terraform {
  required_version = "~> 1.8"

  required_providers {
    utilities = {
      source  = "craigsloggett/utility-functions"
      version = "0.1.0"
    }
  }
}
```

**providers.tf**
```terraform
provider "utilities" {}
```