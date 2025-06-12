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
	_ datasource.DataSource              = &routingTargetDataSource{}
	_ datasource.DataSourceWithConfigure = &routingTargetDataSource{}
)

// NewRoutingTargetDataSource is a helper function to simplify the provider implementation.
func NewRoutingTargetDataSource() datasource.DataSource {
	return &routingTargetDataSource{}
}

// routingTargetDataSource is the data source implementation.
type routingTargetDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *routingTargetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *routingTargetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_targets"
}

// routingTargetDataSourceModel maps the data source schema data.
type routingTargetDataSourceModel struct {
	RoutingTargets []routingTargetReadModel `tfsdk:"routing_targets"`
	Filter         types.Object             `tfsdk:"filter"`
}

// routingTargetReadModel maps schema data.
type routingTargetReadModel struct {
	Id                           types.String        `tfsdk:"id"`
	Name                         types.String        `tfsdk:"name"`
	Moniker                      types.String        `tfsdk:"moniker"`
	RoutingTargetType            routingTargetType   `tfsdk:"routing_target_type"`
	RoutingTargetStatus          routingTargetStatus `tfsdk:"routing_target_status"`
	VSlice                       vSlice              `tfsdk:"vslice"`
	ConfigurationData            types.String        `tfsdk:"configuration_data"`
	PublicInstanceConfiguration  types.String        `tfsdk:"public_instance_configuration"`
	RoutingTargetTypeInstanceId  types.Int32         `tfsdk:"routing_target_type_instance_id"`
	RoutingRedundancyZoneName    types.String        `tfsdk:"routing_redundancy_zone_name"`
	RoutingRedundancyZoneMoniker types.String        `tfsdk:"routing_redundancy_zone_moniker"`
	RegionalGatewayMoniker       types.String        `tfsdk:"regional_gateway_moniker"`
	RegionalGatewayName          types.String        `tfsdk:"regional_gateway_name"`
}

type routingTargetType struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type routingTargetStatus struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type vSlice struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

// Schema defines the schema for the data source.
func (d *routingTargetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of Routing Targets.",
		Attributes: map[string]schema.Attribute{
			"routing_targets": schema.ListNestedAttribute{
				Description: "List of Routing Targets.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Placeholder identifier attribute.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the Routing Target",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the Routing Target",
							Computed:    true,
						},
						"routing_target_type": schema.SingleNestedAttribute{
							Description: "The type of target such as Internet, VPN, Wireguard ",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
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
						"routing_target_status": schema.SingleNestedAttribute{
							Description: "Indicates if it is connected to its target or its current status",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
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
						"vslice": schema.SingleNestedAttribute{
							Description: "The VSlice to link to the routing target",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
							},
						},
						"configuration_data": schema.StringAttribute{
							Description: "JSON config data used for VPNs",
							Computed:    true,
						},
						"public_instance_configuration": schema.StringAttribute{
							Description: "JSON of instance config including the public IP",
							Computed:    true,
						},
						"routing_redundancy_zone_name": schema.StringAttribute{
							Description: "The Redundancy Zone is based on the Region.",
							Computed:    true,
						},
						"routing_redundancy_zone_moniker": schema.StringAttribute{
							Description: "Fallback redudancy zone Moniker",
							Computed:    true,
						},
						"routing_target_type_instance_id": schema.Int32Attribute{
							Computed: true,
						},
						"regional_gateway_moniker": schema.StringAttribute{
							Description: "The Moniker of the region gateway",
							Computed:    true,
						},
						"regional_gateway_name": schema.StringAttribute{
							Description: "The region gateway that you are targeting",
							Computed:    true,
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
func (d *routingTargetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state routingTargetDataSourceModel
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

	routingTargets, err := d.client.GetRoutingTargets(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity Routing Target",
			err.Error(),
		)
		return
	}

	for _, routingTarget := range routingTargets {

		routingTargetState := routingTargetReadModel{}
		err = stacuity.ConvertFromAPI(routingTarget, &routingTargetState)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Stacuity routing targets",
				err.Error(),
			)
			return
		}

		state.RoutingTargets = append(state.RoutingTargets, routingTargetState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
