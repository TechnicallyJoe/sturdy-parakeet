# Demo Storage Account Component
terraform {
  required_version = ">= 1.0.0"
}

provider "azurerm" {
  features {}
}

variable "resource_group_name" {
  type        = string
  description = "The name of the resource group"
}

variable "name" {
  type        = string
  description = "The name of the storage account"
}

variable "location" {
  type        = string
  description = "The Azure region"
  default     = "westeurope"
}



resource "azurerm_storage_account" "main" {
  name                     = var.name
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

output "name" {
  value       = var.name
  description = "The storage account name"
}

output "id" {
  value       = azurerm_storage_account.main.id
  description = "The storage account ID"
}
