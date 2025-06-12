# Copyright (c) HashiCorp, Inc.

resource "stacuity_event_map" "test_event_map" {
  name        = "Terraform Event Map VSlice"
  moniker     = "terraform-eventmap-vslice-comb"
  event_scope = "vslice"
  subscriptions = [
    {
      event_endpoint_id = stacuity_event_handler.test_event_handler.moniker,
      event_type_id     = "vpnchildsaphase2up_v1"
    },
    {
      event_endpoint_id = stacuity_event_handler.test_event_handler.moniker,
      event_type_id     = "vpninitiationsucceeded_v1"
    }
  ]
}