variable "project_name" {
  description = "The name of the project"
  type        = string
}

variable "database_username" {
  description = "Username for the postgres DB user."
  type        = string
}

variable "database_password" {
  description = "password of the postgres database"
  type        = string
}

variable "database_port" {
  description = "port used by postgres database"
  type        = number
}

variable "database_host" {
  default = "Postgres database host's address"
  type    = string
}

variable "database_name" {
  default = "Name of postgres database"
  type    = string
}

variable "image_registry" {
  description = "The registry where the image is stored"
  type        = string
  default     = "ghcr.io/soat1stackgolang"
}

variable "svc_hackaton_image_tag" {
  description = "The tag of the image for the Hackaton service"
  type        = string
  default     = "svc-develop"
}

variable "redeploy_annotation" {
  description = "Annotation to trigger a redeploy"
  type        = string
  default     = "none"
}
