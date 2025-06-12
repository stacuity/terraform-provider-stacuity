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
	_ datasource.DataSource              = &eventMapDataSource{}
	_ datasource.DataSourceWithConfigure = &eventMapDataSource{}
)

// NewEventMapDataSource is a helper function to simplify the provider implementation.
func NewEventMapDataSource() datasource.DataSource {
	return &eventMapDataSource{}
}

// eventMapDataSource is the data source implementation.
type eventMapDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *eventMapDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *eventMapDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_maps"
}

// eventMapsDataSourceModel maps the data source schema data.
type eventMapsDataSourceModel struct {
	EventMaps []eventMapReadModel `tfsdk:"eventmaps"`
	Filter    types.Object        `tfsdk:"filter"`
}

// eventMapReadModel maps schema data.
type eventMapReadModel struct {
	Id            types.String    `tfsdk:"id"`
	Moniker       types.String    `tfsdk:"moniker"`
	Name          types.String    `tfsdk:"name"`
	EventScope    eventScope      `tfsdk:"event_scope"`
	Subscriptions []*subscription `tfsdk:"subscriptions"`
}

type subscription struct {
	EventMap      eventMap      `tfsdk:"event_map"`
	EventEndpoint eventEndpoint `tfsdk:"event_endpoint"`
	EventType     eventType     `tfsdk:"event_type"`
}

type eventEndpoint struct {
	Moniker            types.String `tfsdk:"moniker"`
	Name               types.String `tfsdk:"name"`
	Type               types.String `tfsdk:"type"`
	SummaryDescription types.String `tfsdk:"summary_description"`
}

type eventType struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

type eventScope struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

// Schema defines the schema for the data source.
func (d *eventMapDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of eventMaps.",
		Attributes: map[string]schema.Attribute{
			"eventmaps": schema.ListNestedAttribute{
				Description: "List of event maps.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the event map.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the event map",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the event map",
							Computed:    true,
						},
						"event_scope": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker for the type of event scope",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the event scope",
									Computed:    true,
								},
								"active": schema.BoolAttribute{
									Description: "Active status of the event scope",
									Computed:    true,
								},
							},
						},
						"subscriptions": schema.SetNestedAttribute{
							Description: "List of subscriptions attached to event map.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"event_type": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for the type of event scope",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the event scope",
												Computed:    true,
											},
											"active": schema.BoolAttribute{
												Description: "Active status of the event scope",
												Computed:    true,
											},
										},
									},
									"event_map": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for the type of event scope",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the event scope",
												Computed:    true,
											},
										},
									},
									"event_endpoint": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"moniker": schema.StringAttribute{
												Description: "API Moniker for the type of event scope",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "Name of the event scope",
												Computed:    true,
											},
											"type": schema.StringAttribute{
												Description: "Name of the event handler",
												Computed:    true,
											},
											"summary_description": schema.StringAttribute{
												Description: "Basic description",
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
func (d *eventMapDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state eventMapsDataSourceModel
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

	eventMaps, err := d.client.GetEventMaps(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity event maps",
			err.Error(),
		)
		return
	}
	for _, eventMap := range eventMaps {

		subscription, err := d.client.GetEventMapSubscriptions(eventMap.Moniker)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Stacuity event maps subscriptions",
				err.Error(),
			)
			return
		}

		eventMap.Subscriptions = &subscription

		eventMapState := eventMapReadModel{}
		err = stacuity.ConvertFromAPI(eventMap, &eventMapState)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Stacuity event maps",
				err.Error(),
			)
			return
		}

		state.EventMaps = append(state.EventMaps, eventMapState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
