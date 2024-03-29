variable "region" {
  description = "AWS region to deploy to"
  default     = "us-east-1"
  type        = string
}

variable "availability_zones" {
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
  description = "List of availability zones for the selected region"
}


variable "vpc_cidr_block" {
  type        = string
  default     = "10.0.0.0/16"
  description = "CIDR block range for vpc"
}

variable "project_name" {
  type    = string
  default = "hackaton"
}

variable "database_credentials" {
  description = "Credentials for database creation"

  type = object({
    username = string
    password = string
    port     = string
    name     = string
  })

  default = {
    username = "databaseteste"
    password = "password"
    port     = 5432
    name     = "hackaton"
  }
}

variable "cognito_test_user" {
  description = "Credentials of Cognito user to be used in validation"

  type = object({
    username = string
    password = string
  })

  default = {
    username = "11122233300"
    password = "F@ap1234"
  }
}

variable "redeploy_annotation" {
  description = "Annotation to trigger a redeploy"
  type        = string
  default     = "none"
}

variable "environment" {
  description = "Environment to deploy the stack"
  type        = string
  default     = "dev"
}

variable "svc_hackaton_image_tag" {
  description = "Tag for the image to be used in the deployment"
  type        = string
  default     = "develop"
}