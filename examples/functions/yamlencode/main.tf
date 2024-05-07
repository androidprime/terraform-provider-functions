locals {
  sources = [
    {
      filename = "example/Chart.yaml"
      content = replace(yamlencode({
        apiVersion = "v2"
        name       = "example"
        type       = "application"
        version    = "0.0.2"
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
  base64zip = provider::functions::yamlencode("asdf")
  //base64zip = (
  //  provider::functions::yamlencode({
  //    apiVersion = "v2"
  //    name       = "example"
  //    type       = "application"
  //    version    = "0.0.2"
  //  })
  //)
}
