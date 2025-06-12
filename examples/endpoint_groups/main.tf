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

data "stacuity_endpoint_groups" "endpoint_groups_data" {
  filter = {
    filter  = "name:terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_endpoint_group" "test_endpoint_group_basic" {
  name                    = "terraform endpoint basic"
  moniker                 = "tf-basic-group"
  vslice                  = "tf-test"
  regional_gateway_policy = "automatic"
  ip_allocation_type      = "static"
}

resource "stacuity_endpoint_group" "test_endpoint_group_advanced" {
  name                    = "terraform endpoint advanced"
  moniker                 = "tf-group"
  vslice                  = "tf-test"
  regional_gateway_policy = "automatic"
  ip_allocation_type      = "static"
  event_map               = "vslice-map"
  routing_policy          = "slice1-test-rp"
  operator_policy         = "iomonly"
}
