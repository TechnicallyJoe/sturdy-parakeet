# Demo Key Vault Component
terraform {
  required_version = ">= 1.0.0"
}

variable "name" {
  type        = string
  description = "The name of the key vault"
  default     = "demo-keyvault"
}

variable "resource_group_name" {
  type        = string
  description = "The name of the resource group"
  default     = "demo-rg"
}

variable "location" {
  type        = string
  description = "The Azure region"
  default     = "eastus"
}

variable "sku_name" {
  type        = string
  description = "The SKU of the key vault"
  default     = "standard"
}

output "name" {
  value       = var.name
  description = "The key vault name"
}
