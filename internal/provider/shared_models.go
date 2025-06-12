// Copyright (c) HashiCorp, Inc.

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type filterModel struct {
	Offset types.Int32  `tfsdk:"offset"`
	Limit  types.Int32  `tfsdk:"limit"`
	Filter types.String `tfsdk:"filter"`
	SortBy types.String `tfsdk:"sort_by"`
}
