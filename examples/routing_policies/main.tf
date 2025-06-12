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

data "stacuity_routing_policies" "routing_policies_data" {
  filter = {
    filter  = "name:terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_routing_policy" "test_routing_policy_onerule" {
  name                  = "terraform drop packets"
  moniker               = "tf-drop-packets"
  vslice                = "tf-test"
  routing_policy_status = "active"

  routing_policy_rules = [{
    description            = "terraform drop packets."
    rule_action            = "drop"
    rule_direction         = "uplink"
    destination_ip_pattern = "4.3.2.1/32"
    reflexive              = true
    regional_gateway       = "europe"
    enabled                = true
  }]

  # routing_policy_edge_services = [{
  #   moniker                   = "remoteaccessproxy"
  #   enabled                   = true,
  #   edge_service_instance_ids = ["your_edge_service_instance_moniker"]
  # }]

  rate_limit_uplink_moniker          = "1mbits"
  rate_limit_downlink_moniker        = "1kbits"
  packet_discard_uplink_percentage   = 5
  packet_discard_downlink_percentage = 70
}

resource "stacuity_routing_policy" "test_routing_policy_tworules" {
  name                  = "terraform forward packets"
  moniker               = "terraform-forward-packets"
  vslice                = "tf-test"
  routing_policy_status = "active"

  routing_policy_rules = [
    {
      description              = "terraform forward TCP packets."
      rule_action              = "forward"
      rule_direction           = "uplink"
      destination_ip_pattern   = "4.3.2.1/32"
      divert_ip                = "5.6.7.8"
      divert_port              = "1234"
      transport_protocol       = "tcp"
      source_port_pattern      = "1111"
      destination_port_pattern = "4321"
      routing_target           = "tf-test_routing_target"
      reflexive                = true
      regional_gateway         = "europe"
      enabled                  = true
    },
    {
      description        = "terraform drop UDP packets."
      rule_action        = "drop"
      rule_direction     = "downlink"
      source_ip_pattern  = "1.2.3.4/32"
      divert_ip          = "5.6.7.8"
      divert_port        = "1234"
      transport_protocol = "udp"
      reflexive          = true
      regional_gateway   = "europe"
      enabled            = true
    }
  ]

  # routing_policy_edge_services = [{
  #   moniker                   = "remoteaccessproxy"
  #   enabled                   = true,
  #   edge_service_instance_ids = ["your_edge_service_instance_moniker"]
  # }]

  rate_limit_uplink_moniker          = "10mbits"
  rate_limit_downlink_moniker        = "unlimited"
  packet_discard_uplink_percentage   = 1
  packet_discard_downlink_percentage = 5
}


