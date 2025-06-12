# Copyright (c) HashiCorp, Inc.

output "routing_policies" {
  description = "All routing policies"
  value       = data.stacuity_routing_policies.routing_policies_data
}