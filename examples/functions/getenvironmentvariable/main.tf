locals {
  key = "PING"
  value = provider::functions::getenvironmentvariable(local.key)
}
