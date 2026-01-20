# Demo Kubernetes ArgoCD Base
terraform {
  required_version = ">= 1.0.0"
}

variable "namespace" {
  type        = string
  description = "The Kubernetes namespace"
  default     = "argocd"
}

variable "argocd_version" {
  type        = string
  description = "The ArgoCD version"
  default     = "v2.9.0"
}

variable "enable_ha" {
  type        = bool
  description = "Enable high availability mode"
  default     = false
}

output "namespace" {
  value       = var.namespace
  description = "The ArgoCD namespace"
}

output "version" {
  value       = var.argocd_version
  description = "The ArgoCD version"
}
