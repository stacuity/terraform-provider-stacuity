# Copyright (c) HashiCorp, Inc.

//Inital run returns nothing(if nothing exists) but then second apply will show data as it will now exist
output "vslices_filtered_terraform" {
  value = data.stacuity_vslices.vslice_data
}