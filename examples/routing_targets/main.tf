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

data "stacuity_routing_targets" "routing_targets_data" {
  filter = {
    filter  = "name:Terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_routing_target" "test_routing_target_internet" {
  name                            = "Terraform Internet Target"
  moniker                         = "tf-test_routing_target"
  redundancy_zone_moniker         = "north-america-primary"
  vslice                          = "tf-test" //use one that already exists
  routing_target_type             = "internet"
  routing_target_type_instance_id = "dc14-prod-nat-01a" //internet
}

resource "stacuity_routing_target" "test_routing_target_vpn" {
  name                    = "Terraform VPN Target"
  moniker                 = "tf-vpn"
  redundancy_zone_moniker = "europe-primary"
  configuration_data = {
    vpn_config = {
      remote_peer_address      = "192.168.1.1"
      remote_subnets           = "192.168.0.0/16"
      remote_encryption_domain = "example_domain"
      local_encryption_domain  = "example_local_domain"
      local_subnets            = "10.0.0.0/8"
      preshared_key            = "example_preshared_key"
      key_exchange_type        = "ikev2"
      vpn_ike_option           = "aes128-sha1-modp1024"
      vpn_esp_option           = "aes128-sha1-modp1024"
      phase1_lifetime          = 100
      phase2_lifetime          = 200
    }
  }
  vslice                          = "tf-test" //use one that already exists
  routing_target_type             = "vpn"
  routing_target_type_instance_id = "ma5-prod-vpn-01a-ipsec"
}

resource "stacuity_routing_target" "test_routing_target_wireguard" {
  name                    = "Terraform Wireguard Target"
  moniker                 = "tf-wireguard"
  redundancy_zone_moniker = "europe-primary"
  configuration_data = {
    wireguard_config = {
      local_subnets          = "10.0.0.0/8"
      remote_public_key      = "10.0.0.0/8"
      remote_subnets         = "192.168.0.0/16"
      remote_peer_ip_address = "192.168.1.5"
    }
  }
  vslice                          = "tf-test" //use one that already exists
  routing_target_type             = "wireguard"
  routing_target_type_instance_id = "ma5-prod-vpn-01a-wg"
}

