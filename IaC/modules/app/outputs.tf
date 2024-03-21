
output "lb_service_name_hackaton" {
  value       = local.lb_service_name_hackaton
  description = "Name of the Load Balancer K8s service that exposes Hackaton service"
}

output "lb_service_port_hackaton" {
  value       = local.lb_service_port_hackaton
  description = "Port exposed of the Load Balancer K8s service associated to Hackaton micorservice"
}
