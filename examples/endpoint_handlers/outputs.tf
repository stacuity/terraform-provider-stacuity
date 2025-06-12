# Copyright (c) HashiCorp, Inc.

output "event_handlers" {
  description = "All endpoint handlers"
  value       = data.stacuity_event_handlers.event_handlers_data
}