locals {
  base64targz = (
    provider::functions::base64targz(var.sources)
  )
}
