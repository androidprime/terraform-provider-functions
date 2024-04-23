terraform {
  required_providers {
    hashicups = {
      source  = "registry.terraform.io/hashicorp/functions"
    }
  }
  required_version = ">= 1.8.0"
}

output "total_price" {
  value = provider::hashicups::base64targz([
    {
      filename = "test/foo.txt"
      contents = "foo"
    },
    {
      filename = "test/bar.txt"
      contents = "bar"
    },    
  ])
}
