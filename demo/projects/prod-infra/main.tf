# Demo Production Infrastructure Project
terraform {
  required_version = ">= 1.0.0"
}

variable "environment" {
  type        = string
  description = "The environment name"
  default     = "production"
}

variable "region" {
  type        = string
  description = "The deployment region"
  default     = "eastus"
}

variable "project_name" {
  type        = string
  description = "The project name"
  default     = "prod-infra"
}

output "environment" {
  value       = var.environment
  description = "The environment"
}

output "region" {
  value       = var.region
  description = "The region"
}
