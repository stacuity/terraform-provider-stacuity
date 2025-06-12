# Copyright (c) HashiCorp, Inc.

resource "stacuity_endpoint_group" "test_endpoint_group" {
  name                    = "Terraform Endpoints Combined"
  moniker                 = "terraform-endpoints-comb"
  vslice                  = stacuity_vslice.terraform_combined_vslice.moniker
  routing_policy          = stacuity_routing_policy.test_routing_policy_one_rule.moniker
  regional_gateway_policy = stacuity_regional_policy.test_regional_policy.moniker # or "automatic"
  ip_allocation_type      = "static"
  operator_policy         = stacuity_operator_policy.test_operator_policy.moniker
}