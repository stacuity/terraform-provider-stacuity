# Copyright (c) HashiCorp, Inc.

resource "stacuity_vslice" "terraform_combined_vslice" {
  name               = "Terraform Combined Example"
  moniker            = "terraform-comb"
  subnet_address     = "100.64.0.0/10"
  event_map          = stacuity_event_map.test_event_map.moniker #Optional
  dns_mode           = "custom"
  dns_servers        = ["1.1.1.1", "8.8.8.8"] #Optional if dnsmode = custom
  ip_address_family  = "ipv4"
  ip_allocation_type = "static"
}