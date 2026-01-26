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

variable "account_tier" {
  type        = string
  description = "The storage account tier"
  default     = "Standard"
}

variable "replication_type" {
  type        = string
  description = "The replication type for the storage account"
  default     = "LRS"
}

variable "enable_https_only" {
  type        = bool
  description = "Whether to enforce HTTPS only"
  default     = true
}

variable "container_delete_retention_days" {
  type        = number
  description = "Number of days to retain deleted containers"
  default     = 7
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply to the storage account"
  default     = {}
}

variable "allowed_ips" {
  type        = list(string)
  description = "List of allowed IP addresses"
  default     = []
}

variable "blob_properties" {
  type = object({
    versioning_enabled  = bool
    change_feed_enabled = bool
  })
  description = "Blob storage properties configuration"
  default = {
    versioning_enabled  = false
    change_feed_enabled = false
  }
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
