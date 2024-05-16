locals {
  key = "FOO"
  value = provider::functions::getenvironmentvariable(local.key)
}
