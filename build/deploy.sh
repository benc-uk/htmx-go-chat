#!/bin/bash

location="uksouth"
resourceGroup="temp"
appName="htmx-go-chat"
image="ghcr.io/benc-uk/htmx-go-chat:latest"

# Create a resource group.
az group create --name $resourceGroup --location $location

# Create container app
az containerapp up \
  --name $appName\
  --resource-group temp \
  --location $location \
  --environment 'app-environment' \
  --image $image \
  --target-port 8000 \
  --ingress external 