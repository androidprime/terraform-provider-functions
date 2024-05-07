output "yamlencode" {
  value = local.base64zip
}

output "x" {
  value = (
    yamlencode({
      apiVersion = "v2"
      name       = "example"
      type       = "application"
      version    = "0.0.2"
    })
  )
}