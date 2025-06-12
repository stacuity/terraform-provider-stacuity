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

data "stacuity_regional_policies" "regional_policies_data" {
  filter = {
    filter  = "name:terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_regional_policy" "test_regional_policy" {
  name    = "terraform regional policy"
  moniker = "tf-regional-policy"
  entries = [
    {
      #Default and Mandatory
      regional_gateway_id = "europe"
    },
    {
      iso_3               = "DZA" #Algeria
      regional_gateway_id = "north-america"
    },
    {
      iso_3               = "DZA" #Algeria
      operator_id         = 2145, #Wataniya Telecom
      regional_gateway_id = "north-america"
    },
    {
      iso_3               = "AND" #Andorra
      operator_id         = 122,  #Andorra Telecom
      regional_gateway_id = "europe"
    },
  ]
}