terraform {
  required_providers {
    functions = {
      source = "registry.terraform.io/hashicorp/functions"
    }
  }
  required_version = ">= 1.8.0"
}

output "output" {
  value = (
    provider::functions::base64targz([
      {
        filename = "test/foo.txt"
        contents = "foo"
      },
      {
        filename = "test/bar.txt"
        contents = "bar"
      }
    ])
  )
}

/*
terraform output -json \
| jq .output.value -r \
| base64 -d > output.tar.gz \
| tar xvzf -
*/
