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
	_ datasource.DataSource              = &regionalPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &regionalPolicyDataSource{}
)

// NewRegionalPolicyDataSource is a helper function to simplify the provider implementation.
func NewRegionalPolicyDataSource() datasource.DataSource {
	return &regionalPolicyDataSource{}
}

// RegionalPolicyDataSource is the data source implementation.
type regionalPolicyDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *regionalPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *regionalPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regional_policies"
}

// RegionalPolicysDataSourceModel maps the data source schema data.
type RegionalPolicysDataSourceModel struct {
	RegionalPolicys []RegionalPolicyReadModel `tfsdk:"regionalpolicies"`
	Filter          types.Object              `tfsdk:"filter"`
}

// RegionalPolicyReadModel maps schema data.
type RegionalPolicyReadModel struct {
	Id      types.String           `tfsdk:"id"`
	Moniker types.String           `tfsdk:"moniker"`
	Name    types.String           `tfsdk:"name"`
	IsFixed types.Bool             `tfsdk:"is_fixed"`
	Active  types.Bool             `tfsdk:"active"`
	Entries []*RegionalPolicyEntry `tfsdk:"entries"`
}

type RegionalPolicyEntry struct {
	Id                      types.String    `tfsdk:"id"`
	OperatorId              types.Int32     `tfsdk:"operator_id"`
	Iso3                    types.String    `tfsdk:"iso_3"`
	RegionalGatewayPolicyId types.String    `tfsdk:"regional_gateway_policy_id"`
	RegionalGateway         regionalGateway `tfsdk:"regional_gateway"`
}

type regionalGateway struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
}

// Schema defines the schema for the data source.
func (d *regionalPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of RegionalPolicys.",
		Attributes: map[string]schema.Attribute{
			"regionalpolicies": schema.ListNestedAttribute{
				Description: "List of regional policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the regional policy.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the regional policy",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the regional policy",
							Computed:    true,
						},
						"active": schema.BoolAttribute{
							Description: "If the policy is active",
							Computed:    true,
						},
						"is_fixed": schema.BoolAttribute{
							Description: "If the policy is a default Stacuity policy",
							Computed:    true,
						},
						"entries": schema.SetNestedAttribute{
							Description: "List of entry rules attached to regional policy.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Entry id",
										Computed:    true,
									},
									"operator_id": schema.Int32Attribute{
										Description: "The operator id to apply to",
										Computed:    true,
									},
									"iso_3": schema.StringAttribute{
										Description: "The iso id the rule applies to",
										Computed:    true,
									},
									"regional_gateway_policy_id": schema.StringAttribute{
										Description: "The policy id the rule applies to",
										Computed:    true,
									},
									"regional_gateway": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for the type of gateway",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the gateway",
												Computed:    true,
											},
										},
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
func (d *regionalPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RegionalPolicysDataSourceModel
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

	regionalPolicies, err := d.client.GetRegionalPolicies(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity regional policies",
			err.Error(),
		)
		return
	}
	for _, RegionalPolicy := range regionalPolicies {

		subscription, err := d.client.GetRegionalPolicyEntries(RegionalPolicy.Moniker)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Stacuity regional policy entries",
				err.Error(),
			)
			return
		}

		RegionalPolicy.Entries = &subscription

		RegionalPolicyState := RegionalPolicyReadModel{}
		err = stacuity.ConvertFromAPI(RegionalPolicy, &RegionalPolicyState)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Stacuity regional policies",
				err.Error(),
			)
			return
		}

		state.RegionalPolicys = append(state.RegionalPolicys, RegionalPolicyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
