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
	_ datasource.DataSource              = &routingPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &routingPolicyDataSource{}
)

// NewRoutingPolicyDataSource is a helper function to simplify the provider implementation.
func NewRoutingPolicyDataSource() datasource.DataSource {
	return &routingPolicyDataSource{}
}

// routingPolicyDataSource is the data source implementation.
type routingPolicyDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *routingPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *routingPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_policies"
}

// routingPolicysDataSourceModel maps the data source schema data.
type routingPolicysDataSourceModel struct {
	RoutingPolicies []routingPolicyReadModel `tfsdk:"routingpolicies"`
	Filter          types.Object             `tfsdk:"filter"`
}

// routingPolicyReadModel maps schema data.
type routingPolicyReadModel struct {
	Id                              types.String                     `tfsdk:"id"`
	Name                            types.String                     `tfsdk:"name"`
	Moniker                         types.String                     `tfsdk:"moniker"`
	VSlice                          routingPolicyVSliceModel         `tfsdk:"vslice"`
	RateLimitUplink                 routingPolicyRateLimitModel      `tfsdk:"rate_limit_uplink"`
	RateLimitDownlink               routingPolicyRateLimitModel      `tfsdk:"rate_limit_downlink"`
	PacketDiscardUplinkPercentage   types.Int32                      `tfsdk:"packet_discard_uplink_percentage"`
	PacketDiscardDownlinkPercentage types.Int32                      `tfsdk:"packet_discard_downlink_percentage"`
	RoutingPolicyStatus             routingPolicyStatusModel         `tfsdk:"routing_policy_status"`
	RoutingPolicyRules              []*routingPolicyRuleModel        `tfsdk:"routing_policy_rules"`
	RoutingPolicyEdgeServices       []*routingPolicyEdgeServiceModel `tfsdk:"routing_policy_edge_services"`
}

// Define nested structs for relationships and sub-elements
type routingPolicyVSliceModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type routingPolicyRateLimitModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type routingPolicyStatusModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type routingPolicyRuleModel struct {
	Id                     types.String                         `tfsdk:"id"`
	RoutingPolicyId        types.String                         `tfsdk:"routing_policy_id"`
	Description            types.String                         `tfsdk:"description"`
	RuleAction             routingPolicyRuleActionModel         `tfsdk:"rule_action"`
	RuleDirection          routingPolicyRuleDirectionModel      `tfsdk:"rule_direction"`
	Precedence             types.Int32                          `tfsdk:"precedence"`
	SourceIpPattern        types.String                         `tfsdk:"source_ip_pattern"`
	DestinationIpPattern   types.String                         `tfsdk:"destination_ip_pattern"`
	DivertIp               types.String                         `tfsdk:"divert_ip"`
	DivertPort             types.String                         `tfsdk:"divert_port"`
	TransportProtocol      *routingPolicyTransportProtocolModel `tfsdk:"transport_protocol"`
	SourcePortPattern      types.String                         `tfsdk:"source_port_pattern"`
	DestinationPortPattern types.String                         `tfsdk:"destination_port_pattern"`
	RoutingTarget          routingPolicyRoutingTargetModel      `tfsdk:"routing_target"`
	Reflexive              types.Bool                           `tfsdk:"reflexive"`
	Enabled                types.Bool                           `tfsdk:"enabled"`
	RegionalGateway        routingPolicyRegionalGatewayModel    `tfsdk:"regional_gateway"`
}

type routingPolicyRuleActionModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type routingPolicyRuleDirectionModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type routingPolicyTransportProtocolModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type routingPolicyRoutingTargetModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type routingPolicyRegionalGatewayModel struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

type routingPolicyEdgeServiceModel struct {
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	IconShape              types.String `tfsdk:"icon_shape"`
	Moniker                types.String `tfsdk:"moniker"`
	Available              types.Bool   `tfsdk:"available"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	HasInstance            types.Bool   `tfsdk:"has_instance"`
	EdgeServiceInstanceIds types.String `tfsdk:"edge_service_instance_ids"`
}

// Schema defines the schema for the data source.
func (d *routingPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of routingPolicys.",
		Attributes: map[string]schema.Attribute{
			"routingpolicies": schema.ListNestedAttribute{
				Description: "List of routing policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the routing policy.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the routingPolicy",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the routing policy",
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
						"rate_limit_uplink": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for rate limit uplink",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the rate limit uplink",
									Computed:    true,
								},
								"active": schema.BoolAttribute{
									Description: "Active status of the rate limit uplink",
									Computed:    true,
								},
							},
						},
						"rate_limit_downlink": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for rate limit downlink",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the rate limit downlink",
									Computed:    true,
								},
								"active": schema.BoolAttribute{
									Description: "Active status of the rate limit downlink",
									Computed:    true,
								},
							},
						},
						"packet_discard_uplink_percentage": schema.Int32Attribute{
							Description: "Percentage of packet discard uplink",
							Computed:    true,
						},
						"packet_discard_downlink_percentage": schema.Int32Attribute{
							Description: "Percentage of packet discard downlink",
							Computed:    true,
						},
						"routing_policy_status": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for routing policy status",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the routing policy status",
									Computed:    true,
								},
								"active": schema.BoolAttribute{
									Description: "Active status of the routing policy",
									Computed:    true,
								},
							},
						},
						"routing_policy_rules": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique identifier for the rule",
										Computed:    true,
									},
									"routing_policy_id": schema.StringAttribute{
										Description: "Id of the routing policy to which this rule belongs",
										Computed:    true,
									},
									"description": schema.StringAttribute{
										Description: "Description of the routing rule",
										Computed:    true,
									},
									"rule_action": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for rule action",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the rule action",
												Computed:    true,
											},
											"active": schema.BoolAttribute{
												Description: "Active status of the rule action",
												Computed:    true,
											},
										},
									},
									"rule_direction": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for rule direction",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the rule direction",
												Computed:    true,
											},
											"active": schema.BoolAttribute{
												Description: "Active status of the rule direction",
												Computed:    true,
											},
										},
									},
									"precedence": schema.Int32Attribute{
										Description: "Precedence of the routing rule",
										Computed:    true,
									},
									"source_ip_pattern": schema.StringAttribute{
										Description: "Source IP pattern for the rule",
										Computed:    true,
									},
									"destination_ip_pattern": schema.StringAttribute{
										Description: "Destination IP pattern for the rule",
										Computed:    true,
									},
									"divert_ip": schema.StringAttribute{
										Description: "Divert IP address for the rule",
										Computed:    true,
									},
									"divert_port": schema.StringAttribute{
										Description: "Divert port for the rule",
										Computed:    true,
									},
									"transport_protocol": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for transport protocol",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the transport protocol",
												Computed:    true,
											},
											"active": schema.BoolAttribute{
												Description: "Active status of the transport protocol",
												Computed:    true,
											},
										},
									},
									"source_port_pattern": schema.StringAttribute{
										Description: "Source port pattern for the rule",
										Computed:    true,
									},
									"destination_port_pattern": schema.StringAttribute{
										Description: "Destination port pattern for the rule",
										Computed:    true,
									}, "routing_target": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for routing target",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the routing target",
												Computed:    true,
											},
										},
									},
									"reflexive": schema.BoolAttribute{
										Description: "Reflexive property for the rule",
										Computed:    true,
									},
									"enabled": schema.BoolAttribute{
										Description: "Enabled status of the routing rule",
										Computed:    true,
									},
									"regional_gateway": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for regional gateway",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the regional gateway",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"routing_policy_edge_services": schema.SetNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Name of the edge service",
										Computed:    true,
									},
									"description": schema.StringAttribute{
										Description: "Description of the edge service",
										Computed:    true,
									},
									"icon_shape": schema.StringAttribute{
										Description: "Icon shape for the edge service",
										Computed:    true,
									},
									"moniker": schema.StringAttribute{
										Description: "API Moniker of the edge service",
										Computed:    true,
									},
									"available": schema.BoolAttribute{
										Description: "Availability status of the edge service",
										Computed:    true,
									},
									"enabled": schema.BoolAttribute{
										Description: "Enabled status of the edge service",
										Computed:    true,
									},
									"has_instance": schema.BoolAttribute{
										Description: "Indicates if the edge service has an instance",
										Computed:    true,
									},
									"edge_service_instance_ids": schema.StringAttribute{
										Description: "Unique identifiers for the edge service instance",
										Computed:    true,
									},
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
func (d *routingPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state routingPolicysDataSourceModel
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

	routingPolicies, err := d.client.GetRoutingPolicies(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity routing policies",
			err.Error(),
		)
		return
	}

	for _, routingPolicy := range routingPolicies {
		routingPolicyState := routingPolicyReadModel{}
		err = stacuity.ConvertFromAPI(routingPolicy, &routingPolicyState)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Stacuity routing policies",
				err.Error(),
			)
			return
		}

		state.RoutingPolicies = append(state.RoutingPolicies, routingPolicyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
