// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	stacuity "stacuity.com/go_client"
	models "stacuity.com/go_client/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &vSlicesDataSource{}
	_ datasource.DataSourceWithConfigure = &vSlicesDataSource{}
)

// NewVSliceDataSource is a helper function to simplify the provider implementation.
func NewVSliceDataSource() datasource.DataSource {
	return &vSlicesDataSource{}
}

// vSlicesDataSource is the data source implementation.
type vSlicesDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *vSlicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*stacuity.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *stacuity.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *vSlicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vslices"
}

// vSlicesDataSourceModel maps the data source schema data.
type vSlicesDataSourceModel struct {
	VSlices []vSlicesReadModel `tfsdk:"vslices"`
	Filter  types.Object       `tfsdk:"filter"`
}

// vSlicesReadModel maps vslice schema data.
type vSlicesReadModel struct {
	Id              types.String         `tfsdk:"id"`
	Name            types.String         `tfsdk:"name"`
	Moniker         types.String         `tfsdk:"moniker"`
	Subnets         []types.String       `tfsdk:"subnets"`
	DNSServers      []types.String       `tfsdk:"dns_servers"`
	EventMap        eventMapModel        `tfsdk:"event_map"`
	DNSMode         dnsModeModel         `tfsdk:"dns_mode"`
	IpAddressFamily ipAddressFamilyModel `tfsdk:"ip_address_family"`
}

type eventMapModel struct {
	Id      types.String `tfsdk:"id"`
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type dnsModeModel struct {
	Key     types.Int32  `tfsdk:"key"`
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type ipAddressFamilyModel struct {
	Key     types.Int32  `tfsdk:"key"`
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

// Schema defines the schema for the data source.
func (d *vSlicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of vSlices.",
		Attributes: map[string]schema.Attribute{
			"vslices": schema.ListNestedAttribute{
				Description: "List of vSlices.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Placeholder identifier attribute.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the vSlice",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the vSlice",
							Computed:    true,
						},
						"subnets": schema.ListAttribute{
							Description: "Subnets applied to vSlice.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"dns_servers": schema.ListAttribute{
							Description: "DNS servers applied to vSlice.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"event_map": schema.SingleNestedAttribute{
							Description: "The event map linked to the vSlice",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed: true,
								},
								"moniker": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"dns_mode": schema.SingleNestedAttribute{
							Description: "The DNS mode applied to the vSlice",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"key": schema.Int32Attribute{
									Computed: true,
								},
								"moniker": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"active": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
						"ip_address_family": schema.SingleNestedAttribute{
							Description: "The IP address family type",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"key": schema.Int32Attribute{
									Computed: true,
								},
								"moniker": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"active": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			"filter": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"filter": schema.StringAttribute{
						Description: "Filter the results. Example 'name:TerraForm Test,monkier:tf-test'",
						Optional:    true,
					},
					"sort_by": schema.StringAttribute{
						Description: "Sort by any property. Example 'asc(property),desc(property)'",
						Optional:    true,
					},
					"limit": schema.Int32Attribute{
						Description: "How many results to return",
						Optional:    true,
					},
					"offset": schema.Int32Attribute{
						Description: "What offset to use when querying",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *vSlicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vSlicesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	pagingQuery := models.PagingState{}
	if !state.Filter.IsNull() {
		var filter filterModel
		resp.Diagnostics.Append(state.Filter.As(ctx, &filter, basetypes.ObjectAsOptions{})...)

		pagingQuery = models.PagingState{
			Offset: filter.Offset.ValueInt32(),
			Limit:  filter.Limit.ValueInt32(),
			Filter: filter.Filter.ValueString(),
			SortBy: filter.SortBy.ValueString(),
		}
	}

	vSlices, err := d.client.GetVSlices(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity vSlices",
			err.Error(),
		)
		return
	}

	for _, vSlice := range vSlices {
		vSliceState := vSlicesReadModel{
			Id:      types.StringValue(vSlice.Id),
			Name:    types.StringValue(vSlice.Name),
			Moniker: types.StringValue(vSlice.Moniker),
			DNSMode: dnsModeModel{
				Key:     types.Int32Value(vSlice.DNSMode.Key),
				Moniker: types.StringValue(vSlice.DNSMode.Moniker),
				Active:  types.BoolValue(vSlice.DNSMode.Active),
				Name:    types.StringValue(vSlice.DNSMode.Name),
			},
			EventMap: eventMapModel{
				Id:      types.StringValue(vSlice.EventMap.Id),
				Moniker: types.StringValue(vSlice.EventMap.Moniker),
				Name:    types.StringValue(vSlice.EventMap.Name),
			},
			IpAddressFamily: ipAddressFamilyModel{
				Moniker: types.StringValue(vSlice.IpAddressFamily.Moniker),
				Active:  types.BoolValue(vSlice.IpAddressFamily.Active),
				Name:    types.StringValue(vSlice.IpAddressFamily.Name),
				Key:     types.Int32Value(vSlice.IpAddressFamily.Key),
			},
		}

		for _, subnet := range vSlice.Subnets {
			vSliceState.Subnets = append(vSliceState.Subnets, types.StringValue(subnet))
		}

		for _, dns := range vSlice.DNSServers {
			vSliceState.DNSServers = append(vSliceState.DNSServers, types.StringValue(dns))
		}

		state.VSlices = append(state.VSlices, vSliceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
