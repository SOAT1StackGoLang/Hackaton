data "http" "metrics_server" {
  url = "https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml"
}

resource "local_file" "metrics_server" {
  filename = "${path.module}/metrics_server.yaml"
  content  = data.http.metrics_server.response_body 
}

data "kubectl_file_documents" "metrics_server_documents" {
  content = local_file.metrics_server.content
}

resource "kubectl_manifest" "metrics_server" {
  yaml_body = join("\n---\n", data.kubectl_file_documents.metrics_server_documents.documents)
}