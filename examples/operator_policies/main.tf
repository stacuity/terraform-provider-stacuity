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

data "stacuity_operator_policies" "operator_policies_data" {
  filter = {
    filter  = "name:terraform" #case sensitive contains
    sort_by = "asc(name),desc(moniker)"
    offset  = 0
    limit   = 100
  }
}

resource "stacuity_operator_policy" "test_operator_policy" {
  name    = "terraform operator policy"
  moniker = "tf-operator-policy"
  entries = [
    {
      steering_profile_entry_action = "reject-hard" #Default
    },
    {
      iso_3                         = "AUS" #Australia
      steering_profile_entry_action = "allow"
    },
    {
      iso_3                         = "DZA" #Algeria
      steering_profile_entry_action = "reject-hard"
    },
    {
      iso_3                         = "DZA" #Algeria
      operator_id                   = 2145, #Wataniya Telecom
      steering_profile_entry_action = "allow"
    },
    {
      iso_3                         = "AND" #Andorra
      operator_id                   = 122,  #Andorra Telecom
      steering_profile_entry_action = "reject-soft"
    },
  ]
}