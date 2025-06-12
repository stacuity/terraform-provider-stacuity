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
	_ datasource.DataSource              = &endpointGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &endpointGroupDataSource{}
)

// NewEndpointGroupDataSource is a helper function to simplify the provider implementation.
func NewEndpointGroupDataSource() datasource.DataSource {
	return &endpointGroupDataSource{}
}

// endpointGroupDataSource is the data source implementation.
type endpointGroupDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *endpointGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *endpointGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_groups"
}

// endpointGroupsDataSourceModel maps the data source schema data.
type endpointGroupsDataSourceModel struct {
	EndpointGroups []endpointGroupReadModel `tfsdk:"endpointgroups"`
	Filter         types.Object             `tfsdk:"filter"`
}

// endpointGroupReadModel maps schema data.
type endpointGroupReadModel struct {
	Id                    types.String          `tfsdk:"id"`
	Moniker               types.String          `tfsdk:"moniker"`
	Name                  types.String          `tfsdk:"name"`
	EndpointsAssigned     types.Int32           `tfsdk:"endpoints_assigned"`
	VSlice                vSlice                `tfsdk:"vslice"`
	EventMap              *eventMap             `tfsdk:"event_map"`
	RoutingPolicy         *routingPolicy        `tfsdk:"routing_policy"`
	SteeringProfile       *steeringProfile      `tfsdk:"steering_profile"`
	RegionalGatewayPolicy regionalGatewayPolicy `tfsdk:"regional_gateway_policy"`
	IPAllocationType      ipAllocationType      `tfsdk:"ip_allocation_type"`
	CustomerId            types.String          `tfsdk:"customer_id"`
}

type eventMap struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type routingPolicy struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type steeringProfile struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type regionalGatewayPolicy struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type ipAllocationType struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

// Schema defines the schema for the data source.
func (d *endpointGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of endpointGroups.",
		Attributes: map[string]schema.Attribute{
			"endpointgroups": schema.ListNestedAttribute{
				Description: "List of endpoint groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the endpoint group.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the endpoint group",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the endpoint group",
							Computed:    true,
						},
						"endpoints_assigned": schema.Int32Attribute{
							Description: "Total number of endpoints in the group",
							Computed:    true,
						},
						"customer_id": schema.StringAttribute{
							Description: "Customerid of the endpoint group",
							Computed:    true,
						},
						"vslice": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker of the vslice",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the vslice",
									Computed:    true,
								},
							},
						},
						"event_map": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker of the event map",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the event map",
									Computed:    true,
								},
							},
						},
						"routing_policy": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for routing policy",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the routing policy",
									Computed:    true,
								},
							},
						},
						"steering_profile": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for endpoint group status",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the endpoint group status",
									Computed:    true,
								},
							},
						},

						"regional_gateway_policy": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for endpoint group status",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the endpoint group status",
									Computed:    true,
								},
							},
						},
						"ip_allocation_type": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for the type of ip allocation",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the ip alocation",
									Computed:    true,
								},
								"active": schema.BoolAttribute{
									Description: "Active status of the ip allocation",
									Computed:    true,
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
func (d *endpointGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state endpointGroupsDataSourceModel
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

	routingPolicies, err := d.client.GetEndpointGroups(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity endpoint groups",
			err.Error(),
		)
		return
	}
	for _, endpointGroup := range routingPolicies {
		endpointGroupState := endpointGroupReadModel{}
		err = stacuity.ConvertFromAPI(endpointGroup, &endpointGroupState)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Stacuity endpoint groups",
				err.Error(),
			)
			return
		}

		state.EndpointGroups = append(state.EndpointGroups, endpointGroupState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
