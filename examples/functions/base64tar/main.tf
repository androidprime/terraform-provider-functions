locals {
  sources = [
    {
      filename = "example/Chart.yaml"
      content = replace(yamlencode({
        apiVersion = "v2"
        name       = "example"
        type       = "application"
        version    = "0.0.1"
      }), "\"", "")
    },
    {
      filename = "example/templates/deployment.yaml"
      content = replace(yamlencode({
        apiVersion = "apps/v1"
        kind       = "Deployment"
        metadata = {
          name = "example"
        }
        spec = {
          replicas = 2
          selector = {
            matchLabels = {
              app = "example"
            }
          }
          templates = {
            metadata = {
              labels = {
                app = "example"
              }
              spec = {
                containers = [
                  {
                    name  = "nginx"
                    image = "nginx:latest"
                    ports = [
                      {
                        containerPort = 80
                      }
                    ]
                  }
                ]
              }
            }
          }
        }
      }), "\"", "")
    },
    {
      filename = "example/templates/services.yaml"
      content = replace(yamlencode({
        apiVersion = "v1"
        kind       = "Service"
        metadata = {
          name = "example"
        }
        spec = {
          selector = {
            app = "example"
          }
          ports = [
            {
              port       = 80
              targetPort = 80
            }
          ]
        }
      }), "\"", "")
    }
  ]
  base64tar = (
    provider::functions::base64tar(local.sources)
  )
  base64tar_decoded = base64decode(local.base64tar)
  base64tar_gzip    = base64gzip(local.base64tar_decoded)
}
