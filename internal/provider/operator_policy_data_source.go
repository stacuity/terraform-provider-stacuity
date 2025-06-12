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
	_ datasource.DataSource              = &operatorPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &operatorPolicyDataSource{}
)

// NewOperatorPolicyDataSource is a helper function to simplify the provider implementation.
func NewOperatorPolicyDataSource() datasource.DataSource {
	return &operatorPolicyDataSource{}
}

// OperatorPolicyDataSource is the data source implementation.
type operatorPolicyDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *operatorPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *operatorPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_operator_policies"
}

// OperatorPolicysDataSourceModel maps the data source schema data.
type OperatorPolicysDataSourceModel struct {
	OperatorPolicys []OperatorPolicyReadModel `tfsdk:"operatorpolicies"`
	Filter          types.Object              `tfsdk:"filter"`
}

// OperatorPolicyReadModel maps schema data.
type OperatorPolicyReadModel struct {
	Id       types.String           `tfsdk:"id"`
	Moniker  types.String           `tfsdk:"moniker"`
	Name     types.String           `tfsdk:"name"`
	Allow2g  types.Bool             `tfsdk:"allow_2g"`
	Allow3g  types.Bool             `tfsdk:"allow_3g"`
	Allow45g types.Bool             `tfsdk:"allow_45g"`
	Entries  []*OperatorPolicyEntry `tfsdk:"entries"`
}

type OperatorPolicyEntry struct {
	Id                         types.String               `tfsdk:"id"`
	OperatorId                 types.Int32                `tfsdk:"operator_id"`
	Iso3                       types.String               `tfsdk:"iso_3"`
	SteeringProfileEntryAction SteeringProfileEntryAction `tfsdk:"steering_profile_entry_action"`
}

type SteeringProfileEntryAction struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

// Schema defines the schema for the data source.
func (d *operatorPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of OperatorPolicys.",
		Attributes: map[string]schema.Attribute{
			"operatorpolicies": schema.ListNestedAttribute{
				Description: "List of operator policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the operator policy.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the operator policy",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the operator policy",
							Computed:    true,
						},
						"allow_2g": schema.BoolAttribute{
							Description: "If the policy supports 2g",
							Computed:    true,
						},
						"allow_3g": schema.BoolAttribute{
							Description: "If the policy supports 3g",
							Computed:    true,
						},
						"allow_45g": schema.BoolAttribute{
							Description: "If the policy supports 4g",
							Computed:    true,
						},
						"entries": schema.SetNestedAttribute{
							Description: "List of entry rules attached to operator policy.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "entry id",
										Computed:    true,
									},
									"operator_id": schema.Int32Attribute{
										Description: "the operator id to apply to",
										Computed:    true,
									},
									"iso_3": schema.StringAttribute{
										Description: "the iso id the rule applies to",
										Computed:    true,
									},
									"steering_profile_entry_action": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for the type of steering profile",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the steering profile",
												Computed:    true,
											},
											"active": schema.BoolAttribute{
												Description: "Active status of the steering profile",
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
func (d *operatorPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state OperatorPolicysDataSourceModel
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

	operatorPolicies, err := d.client.GetOperatorPolicies(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity operator policies",
			err.Error(),
		)
		return
	}
	for _, OperatorPolicy := range operatorPolicies {

		subscription, err := d.client.GetOperatorPolicyEntries(OperatorPolicy.Moniker)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Stacuity operator policy entries",
				err.Error(),
			)
			return
		}

		OperatorPolicy.Entries = &subscription

		OperatorPolicyState := OperatorPolicyReadModel{}
		err = stacuity.ConvertFromAPI(OperatorPolicy, &OperatorPolicyState)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Stacuity operator policies",
				err.Error(),
			)
			return
		}

		state.OperatorPolicys = append(state.OperatorPolicys, OperatorPolicyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
