# Copyright (c) HashiCorp, Inc.

resource "stacuity_routing_policy" "test_routing_policy_one_rule" {
  name                  = "terraform drop packets"
  moniker               = "terraform-drop-packets-comb"
  vslice                = stacuity_vslice.terraform_combined_vslice.moniker
  routing_policy_status = "active"

  routing_policy_rules = [{
    description      = "terraform drop packets."
    rule_action      = "drop"
    rule_direction   = "uplink"
    reflexive        = true
    regional_gateway = "europe"
    enabled          = true
  }]

  # routing_policy_edge_services = [{
  #   moniker                   = "remoteaccessproxy"
  #   enabled                   = true,
  #   edge_service_instance_ids = ["5ab67663-5f4a-4b4d-8d22-7cf545392e6f"]
  # }]

  rate_limit_uplink_moniker          = "1mbits"
  rate_limit_downlink_moniker        = "1kbits"
  packet_discard_uplink_percentage   = 5
  packet_discard_downlink_percentage = 70
}

resource "stacuity_routing_policy" "test_routing_policy_multiple_rules" {
  name                  = "terraform forward packets"
  moniker               = "terraform-multiple-rules-comb"
  vslice                = stacuity_vslice.terraform_combined_vslice.moniker
  routing_policy_status = "active"

  routing_policy_rules = [
    {
      description            = "Allow Google Pings"
      rule_action            = "forward"
      rule_direction         = "uplink"
      destination_ip_pattern = "8.8.8.8"
      transport_protocol     = "icmp"
      routing_target         = stacuity_routing_target.test_routing_target_internet.moniker
      reflexive              = true
      enabled                = true
    },
    {
      description              = "Forward TCP packets."
      rule_action              = "forward"
      rule_direction           = "uplink"
      destination_ip_pattern   = "4.3.2.1/32"
      divert_ip                = "5.6.7.8"
      divert_port              = "1234"
      transport_protocol       = "tcp"
      source_port_pattern      = "1111"
      destination_port_pattern = "4321"
      routing_target           = stacuity_routing_target.test_routing_target_wireguard.moniker
      reflexive                = true
      regional_gateway         = "north-america"
      enabled                  = true
    },
    {
      description        = "Drop UDP packets."
      rule_action        = "drop"
      rule_direction     = "downlink"
      source_ip_pattern  = "1.2.3.4/32"
      transport_protocol = "udp"
      reflexive          = true
      regional_gateway   = "europe"
      enabled            = false
    }
  ]

  # routing_policy_edge_services = [{
  #   moniker                   = "remoteaccessproxy"
  #   enabled                   = true,
  #   edge_service_instance_ids = ["5ab67663-5f4a-4b4d-8d22-7cf545392e6f"]
  # }]

  rate_limit_uplink_moniker          = "10mbits"
  rate_limit_downlink_moniker        = "unlimited"
  packet_discard_uplink_percentage   = 1
  packet_discard_downlink_percentage = 5
}