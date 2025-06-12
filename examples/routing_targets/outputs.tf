# Copyright (c) HashiCorp, Inc.

output "routing_targets" {
  description = "All routing targets"
  value       = data.stacuity_routing_targets.routing_targets_data
}