# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    stacuity = {
      source = "registry.terraform.io/stacuity/stacuity"
    }
  }
}

provider "stacuity" {
}

data "stacuity_event_handlers" "event_handlers_data" {
  filter = {
    filter  = "name:terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_event_handler" "test_event_handler_basic" {
  name    = "terraform webhook"
  moniker = "tf-webhook"
  configuration_data = {
    webhook_config = {
      url          = "https://terraform.com",
      timeout      = "10",
      username     = "my username",
      password     = "my password",
      bearer_token = "token"
    }
  }
  event_endpoint_type = "webhook"
}
