terraform {
  required_version = ">= 1.5.7"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm" # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs
      version = ">= 4.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}

module "resource_group" {
  source = "../../"

  name     = "rg-demo-example"
  location = "westeurope"

  tags = {
    Environment = "demo"
    ManagedBy   = "terraform"
  }
}

output "resource_group_name" {
  value = module.resource_group.name
}

output "resource_group_id" {
  value = module.resource_group.id
}
