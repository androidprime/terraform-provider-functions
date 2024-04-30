locals {
  base64tar = (
    provider::functions::base64tar(var.sources)
  )
  //tar = base64decode(local.base64tar)
}
