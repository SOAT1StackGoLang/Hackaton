resource "null_resource" "apply_metrics_server" {
  triggers = {
    metrics_server_yaml = filemd5("${path.module}/metrics_server.yaml")
  }

  provisioner "local-exec" {
    command = "kubectl apply -f ${path.module}/metrics_server.yaml"
  }
}