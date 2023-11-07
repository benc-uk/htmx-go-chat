#!/bin/bash
set -e
red="\e[31m"
nrm="\e[0m"

# =========================================================================
# Deploy with Azure Container App
# This is a simple way to deploy a container to Azure.
# Script is NOT comprehensive and just supports simple scenarios.
# =========================================================================

location="uksouth"
resourceGroup="htmx-chat"
appName="htmx-go-chat"
image="$IMAGE_NAME:$VERSION"

# Check Azure CLI things
echo -e "### Checking if we are logged in to Azure..."
which az >/dev/null || {
  echo -e "$red### Azure CLI is not installed. Goodbye!$nrm"
  exit 1
}
az account show >/dev/null || {
  echo -e "$red### You are not logged in to Azure. Goodbye!$nrm"
  exit 1
}

echo "### Deploying '$image' to Azure"

# Create a resource group.
echo -e "### Creating resource group '$resourceGroup' in '$location'..."
az group create --name "$resourceGroup" --location "$location" >/dev/null

# Create the Azure Container App
echo "### Creating/updating Azure Container App '$appName'..."
az containerapp up \
  --name "$appName" \
  --resource-group "$resourceGroup" \
  --environment 'app-environment' \
  --image "$image" \
  --target-port 8000 \
  --ingress external \
  --browse >/dev/null

# Update the Azure Container App with smaller CPU and memory
echo "### Setting memory and CPU for '$appName'..."
az containerapp update \
  --name "$appName" \
  --resource-group "$resourceGroup" \
  --cpu 0.25 \
  --memory 0.5 >/dev/null
