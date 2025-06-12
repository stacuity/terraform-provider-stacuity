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

data "stacuity_vslices" "vslice_data" {
  filter = {
    filter  = "name:Terraform,moniker:terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_vslice" "test_vslice" {
  name               = "Terraform vSlice"
  moniker            = "terraform-test"
  subnet_address     = "100.64.0.0/10"
  dns_mode           = "auto"
  ip_address_family  = "ipv4"
  ip_allocation_type = "static"
}

resource "stacuity_vslice" "test_vslice2" {
  name               = "Terraform vSlice two"
  moniker            = "terraform-test4"
  subnet_address     = "100.64.0.0/10"
  event_map          = "test1" #Optional
  dns_mode           = "custom"
  dns_servers        = ["1.1.1.1", "8.8.8.8"] #Optional if dnsmode = custom
  ip_address_family  = "ipv4"
  ip_allocation_type = "static"
}