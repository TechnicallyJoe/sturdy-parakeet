# Demo Polylith Repository

This directory contains a demo polylith-style Terraform repository structure for testing and demonstration purposes.

## Structure

```
demo/
├── .tfpl.yml              # tfpl configuration
├── components/            # Reusable terraform components
│   └── azurerm/
│       ├── storage-account/
│       └── key-vault/
├── bases/                 # Composable base configurations
│   └── k8s-argocd/
└── projects/              # Deployable projects
    └── prod-infra/
```

## Usage

From the `demo` directory, you can run tfpl commands:

```bash
cd demo

# Format a component
tfpl fmt -c storage-account

# Initialize a base
tfpl init -b k8s-argocd

# Validate a project (with init)
tfpl val -i -p prod-infra

# Use explicit path
tfpl fmt --path components/azurerm/key-vault

# Pass extra arguments to terraform
tfpl init -c storage-account -a -upgrade
```

## Components

- **storage-account**: Azure Storage Account configuration
- **key-vault**: Azure Key Vault configuration

## Bases

- **k8s-argocd**: Kubernetes ArgoCD deployment base

## Projects

- **prod-infra**: Production infrastructure project
