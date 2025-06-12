# Copyright (c) HashiCorp, Inc.

output "operator_policies" {
  description = "All operator policies"
  value       = data.stacuity_operator_policies.operator_policies_data
}