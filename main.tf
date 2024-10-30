terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.4.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "2.6.0"
    }
  }
}

provider "azurerm" {
  subscription_id = "258b0030-3486-42c7-87c7-2ace4e9e9552"
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

data "archive_file" "function" {
  type        = "zip"
  source_dir  = "${path.module}/build/function"
  output_path = "${path.module}/build/function.zip"
}

resource "azurerm_storage_container" "storage_container_function" {
  name                 = "function-releases"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_storage_blob" "storage_blob_function" {
  name                   = "function-${substr(data.archive_file.function.output_md5, 0, 6)}.zip"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.storage_container_function.name
  type                   = "Block"
  content_md5            = data.archive_file.function.output_md5
  source                 = "${path.module}/build/function.zip"
}

data "archive_file" "app" {
  type        = "zip"
  source_dir  = "${path.module}/build/function"
  output_path = "${path.module}/build/function-${data.archive_file.function.output_sha256}.zip"
}

resource "azurerm_resource_group" "example" {
  name     = "dogop-serverless"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "dogopstorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}


resource "azurerm_log_analytics_workspace" "logs" {
  name                = "dogop-log-analytics"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30

}

resource "azurerm_application_insights" "insights" {
  name                = "dogop-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  workspace_id        = azurerm_log_analytics_workspace.logs.id
  application_type    = "other"
}

#resource "azurerm_app_service_plan" "example" {
#  name                = "azure-functions-test-service-plan"
#  location            = azurerm_resource_group.example.location
#  resource_group_name = azurerm_resource_group.example.name
#  kind                = "Linux"
#  reserved = true

#sku {
#  tier = "Dynamic"
#  size = "Y1"
# }
#}

resource "azurerm_service_plan" "example" {
  name                = "dogop-service-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

resource "azurerm_linux_function_app" "example" {
  name                        = "dogop-serverless-quote"
  resource_group_name         = azurerm_resource_group.example.name
  storage_account_name        = azurerm_storage_account.example.name
  storage_account_access_key  = azurerm_storage_account.example.primary_access_key
  functions_extension_version = "~4"
  location                    = azurerm_resource_group.example.location
  service_plan_id             = azurerm_service_plan.example.id

  zip_deploy_file = data.archive_file.app.output_path
  site_config {
    #application_stack {
    #  use_custom_runtime = true
    #}
    application_insights_connection_string = azurerm_application_insights.insights.connection_string
    application_insights_key               = azurerm_application_insights.insights.instrumentation_key
  }

  app_settings = {
    "WEBSITE_RUN_FROM_PACKAGE" = 1
    #FUNCTIONS_WORKER_RUNTIME=custom
  }

}

output "function_app_name" {
  value       = azurerm_linux_function_app.example.name
  description = "Deployed function app name"
}

output "function_app_default_hostname" {
  value       = azurerm_linux_function_app.example.default_hostname
  description = "Deployed function app hostname"
}
