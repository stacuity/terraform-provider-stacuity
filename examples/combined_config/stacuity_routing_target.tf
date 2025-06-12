# Copyright (c) HashiCorp, Inc.

resource "stacuity_routing_target" "test_routing_target_wireguard" {
  name                    = "Terraform Wireguard Target - Combined"
  moniker                 = "terraform-wireguard-comb"
  redundancy_zone_moniker = "europe-primary"
  configuration_data = {
    wireguard_config = {
      local_subnets          = "10.0.0.0/8"
      remote_public_key      = "10.0.0.0/8"
      remote_subnets         = "192.168.0.0/16"
      remote_peer_ip_address = "192.168.1.5"
    }
  }
  vslice                          = stacuity_vslice.terraform_combined_vslice.moniker
  routing_target_type             = "wireguard"
  routing_target_type_instance_id = "ma5-prod-vpn-01a-wg"
}

resource "stacuity_routing_target" "test_routing_target_internet" {
  name                            = "Terraform Internet Target - Combined"
  moniker                         = "terraform-internet-comb"
  redundancy_zone_moniker         = "north-america-primary"
  vslice                          = stacuity_vslice.terraform_combined_vslice.moniker
  routing_target_type             = "internet"
  routing_target_type_instance_id = "dc14-prod-nat-01a"
}
