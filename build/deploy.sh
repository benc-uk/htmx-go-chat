#!/bin/bash
set -e

# ==========================================================
# Deploy with Azure Container App
# ==========================================================

location="uksouth"
resourceGroup="htmx-chat"
appName="htmx-go-chat"
image="$IMAGE_NAME:$VERSION"

# Create a resource group.
echo -e "\n### Creating resource group '$resourceGroup' in '$location'..."
az group create --name "$resourceGroup" --location "$location" > /dev/null

# Create the Azure Container App
echo "### Creating/updating Azure Container App '$appName'..."
az containerapp up \
  --name "$appName" \
  --resource-group "$resourceGroup" \
  --environment 'app-environment' \
  --image "$image" \
  --target-port 8000 \
  --ingress external \
  --browse > /dev/null

# Update the Azure Container App with smaller CPU and memory
echo "### Setting memory and CPU for '$appName'..."
az containerapp update \
  --name "$appName" \
  --resource-group "$resourceGroup" \
  --cpu 0.25 \
  --memory 0.5 > /dev/null
