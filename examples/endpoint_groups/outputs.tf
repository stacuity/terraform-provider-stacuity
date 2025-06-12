# Copyright (c) HashiCorp, Inc.

output "endpoint_groups" {
  description = "All endpoint groups"
  value       = data.stacuity_endpoint_groups.endpoint_groups_data
}