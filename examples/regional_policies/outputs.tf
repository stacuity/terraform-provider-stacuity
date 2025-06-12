# Copyright (c) HashiCorp, Inc.

output "regional_policies" {
  description = "All regional policies"
  value       = data.stacuity_regional_policies.regional_policies_data
}