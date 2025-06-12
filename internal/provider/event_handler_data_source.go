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
	_ datasource.DataSource              = &eventHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &eventHandlerDataSource{}
)

// NewEventHandlerDataSource is a helper function to simplify the provider implementation.
func NewEventHandlerDataSource() datasource.DataSource {
	return &eventHandlerDataSource{}
}

// eventHandlerDataSource is the data source implementation.
type eventHandlerDataSource struct {
	client *stacuity.Client
}

// Configure adds the provider configured client to the data source.
func (d *eventHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *eventHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_handlers"
}

// eventHandlersDataSourceModel maps the data source schema data.
type eventHandlersDataSourceModel struct {
	EventHandlers []eventHandlerReadModel `tfsdk:"eventhandlers"`
	Filter        types.Object            `tfsdk:"filter"`
}

// eventHandlerReadModel maps schema data.
type eventHandlerReadModel struct {
	Id                 types.String      `tfsdk:"id"`
	Moniker            types.String      `tfsdk:"moniker"`
	Name               types.String      `tfsdk:"name"`
	ConfigurationData  configurationData `tfsdk:"configuration_data"`
	EventEndpointType  eventEndpointType `tfsdk:"event_endpoint_type"`
	SummaryDescription types.String      `tfsdk:"summary_description"`
}

type configurationData struct {
	WebhookConfig *webhookConfig `tfsdk:"webhook_config"`
}

type webhookConfig struct {
	Url         types.String `tfsdk:"url"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Timeout     types.String `tfsdk:"timeout"`
	BearerToken types.String `tfsdk:"bearer_token"`
}

type eventEndpointType struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

// Schema defines the schema for the data source.
func (d *eventHandlerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of eventHandlers.",
		Attributes: map[string]schema.Attribute{
			"eventhandlers": schema.ListNestedAttribute{
				Description: "List of event handlers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier for the event handler.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the event handler",
							Computed:    true,
						},
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the event handler",
							Computed:    true,
						},
						"summary_description": schema.StringAttribute{
							Description: "Summary description for the event handler",
							Computed:    true,
						},
						"configuration_data": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"webhook_config": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"bearer_token": schema.StringAttribute{
											Description: "Bearer token for the webhook",
											Computed:    true,
										},
										"url": schema.StringAttribute{
											Description: "URL for the webhook",
											Computed:    true,
										},
										"username": schema.StringAttribute{
											Description: "Username for the webhook",
											Computed:    true,
										},
										"password": schema.StringAttribute{
											Description: "Password for the webhook",
											Computed:    true,
										},
										"timeout": schema.StringAttribute{
											Description: "Timeout for the webhook",
											Computed:    true,
										},
									},
								},
							},
						},
						"event_endpoint_type": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"moniker": schema.StringAttribute{
									Description: "API Moniker of the event handler",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "Name of the event handler",
									Computed:    true,
								},
								"active": schema.BoolAttribute{
									Description: "Whether the event handler is active",
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
func (d *eventHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state eventHandlersDataSourceModel
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

	eventHandlers, err := d.client.GetEventHandlers(pagingQuery)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Stacuity event handlers",
			err.Error(),
		)
		return
	}
	for _, eventHandler := range eventHandlers {
		eventHandlerState := eventHandlerReadModel{}
		err = stacuity.ConvertFromAPI(eventHandler, &eventHandlerState)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Convert Stacuity event handlers",
				err.Error(),
			)
			return
		}

		state.EventHandlers = append(state.EventHandlers, eventHandlerState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
