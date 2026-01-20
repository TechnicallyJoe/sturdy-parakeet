# Demo Resource Group Component
terraform {
  required_version = ">= 1.0.0"
}

provider "azurerm" {
  features {}
}

variable "name" {
  type        = string
  description = "The name of the resource group"
}

variable "location" {
  type        = string
  description = "The Azure region"
  default     = "westeurope"
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to the resource group"
  default     = {}
}

resource "azurerm_resource_group" "main" {
  name     = var.name
  location = var.location
  tags     = var.tags
}

output "name" {
  value       = azurerm_resource_group.main.name
  description = "The resource group name"
}

output "id" {
  value       = azurerm_resource_group.main.id
  description = "The resource group ID"
}

output "location" {
  value       = azurerm_resource_group.main.location
  description = "The resource group location"
}
