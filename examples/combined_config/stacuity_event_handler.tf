# Copyright (c) HashiCorp, Inc.

resource "stacuity_event_handler" "test_event_handler" {
  name                = "Terraform webhook"
  moniker             = "terraform-webhook-comb"
  event_endpoint_type = "webhook"
  configuration_data = {
    webhook_config = {
      url          = "https://stacuity.com",
      timeout      = "10",
      username     = "Username",
      password     = "P@ssword",
      bearer_token = "T0ken!"
    }
  }
}