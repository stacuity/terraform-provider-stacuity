# Copyright (c) HashiCorp, Inc.

resource "stacuity_regional_policy" "test_regional_policy" {
  name    = "terraform regional policy"
  moniker = "terraform-regional-policy-comb"
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