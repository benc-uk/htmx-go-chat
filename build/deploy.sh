#!/bin/bash

# ==========================================================
# Deploy with Azure Container App
# ==========================================================

location="uksouth"
resourceGroup="htmx-chat"
appName="htmx-go-chat"
image="$IMAGE_NAME:$VERSION"

# Create a resource group.
az group create --name "$resourceGroup" --location "$location"

# Create the Azure Container App
az containerapp up \
  --name "$appName" \
  --resource-group "$resourceGroup" \
  --environment 'app-environment' \
  --image "$image" \
  --target-port 8000 \
  --ingress external \
  --browse 

# Update the Azure Container App with smaller CPU and memory
az containerapp update \
  --name "$appName" \
  --resource-group "$resourceGroup" \
  --cpu 0.25 \
  --memory 0.5 
