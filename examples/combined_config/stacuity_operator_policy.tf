# Copyright (c) HashiCorp, Inc.

resource "stacuity_operator_policy" "test_operator_policy" {
  name    = "terraform custom operator policy"
  moniker = "terraform-comb-operator-policy"
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