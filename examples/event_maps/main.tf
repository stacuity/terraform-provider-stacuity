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

data "stacuity_event_maps" "event_maps_data" {
  filter = {
    filter  = "name:terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_event_map" "test_event_map_basic" {
  name        = "terraform event map"
  moniker     = "tf-event-map"
  event_scope = "vslice"
  subscriptions = [
    {
      event_endpoint_id = "tf-webhook",
      event_type_id     = "vpnchildsaphase2up_v1"
    }
  ]
}

resource "stacuity_event_map" "test_event_map_basic_2" {
  name        = "terraform event map 2"
  moniker     = "tf-event-map-2"
  event_scope = "vslice"
}
