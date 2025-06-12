# Copyright (c) HashiCorp, Inc.

output "event_maps" {
  description = "All event maps"
  value       = data.stacuity_event_maps.event_maps_data
}