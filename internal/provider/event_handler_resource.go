// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	stacuity "stacuity.com/go_client"
	models "stacuity.com/go_client/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &eventHandlerResource{}
	_ resource.ResourceWithConfigure   = &eventHandlerResource{}
	_ resource.ResourceWithImportState = &eventHandlerResource{}
)

// NewEventHandlerResource is a helper function to simplify the provider implementation.
func NewEventHandlerResource() resource.Resource {
	return &eventHandlerResource{}
}

// eventHandlerResource is the resource implementation.
type eventHandlerResource struct {
	client *stacuity.Client
}

type eventHandlerResourceModel struct {
	Id                types.String      `tfsdk:"id"`
	Name              types.String      `tfsdk:"name"`
	Moniker           types.String      `tfsdk:"moniker"`
	ConfigurationData configurationData `tfsdk:"configuration_data"`
	EventEndpointType types.String      `tfsdk:"event_endpoint_type"`
}

func (r *eventHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *eventHandlerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_handler"
}

func (r *eventHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "event handler resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the Event Handler.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the Event Handler.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the Event Handler.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"event_endpoint_type": schema.StringAttribute{
				Description: "The type of event handler. eg webhook",
				Required:    true,
			},
			"configuration_data": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"webhook_config": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"bearer_token": schema.StringAttribute{
								Description: "Bearer token for the webhook",
								Optional:    true,
							},
							"url": schema.StringAttribute{
								Description: "URL for the webhook",
								Required:    true,
							},
							"username": schema.StringAttribute{
								Description: "Username for the webhook",
								Optional:    true,
							},
							"password": schema.StringAttribute{
								Description: "Password for the webhook",
								Optional:    true,
							},
							"timeout": schema.StringAttribute{
								Description: "Timeout for the webhook",
								Required:    true,
							},
						},
					},
				},
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *eventHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *eventHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan eventHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiData := models.EventHandlerModifyItem{}
	var err = stacuity.ConvertToAPI(plan, &apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating event handler",
			"Could not create API event handler, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new event handler
	createResponse, err := r.client.CreateEventHandler(apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating event handler",
			"Could not create event handler, unexpected error: "+err.Error(),
		)
		return
	}

	if createResponse.Success {
		getResponse, err := r.client.GetEventHandler(createResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading event handler",
				"Could not read event handler, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel := eventHandlerResourceModel{}
		err = stacuity.ConvertFromAPI(getResponse, &configDataModel)
		configDataModel.EventEndpointType = types.StringValue(getResponse.EventEndpointType.Moniker)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting to TF",
				"Could not convert from API event handler, unexpected error: "+err.Error(),
			)
			return
		}

		plan = configDataModel
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *eventHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state eventHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed event handler values
	apiResponse, err := r.client.GetEventHandler(state.Moniker.ValueString())

	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading event handler Info",
			"Could not read event handler Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	configDataModel := eventHandlerResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert event handler, unexpected error: "+err.Error(),
		)
		return
	}

	configDataModel.EventEndpointType = types.StringValue(apiResponse.EventEndpointType.Moniker)
	state = configDataModel

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *eventHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan eventHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiConfigDataModel := models.EventHandlerModifyItem{}
	err := stacuity.ConvertToAPI(plan, &apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert event handler, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing event handler
	_, err = r.client.UpdateEventHandler(plan.Moniker.ValueString(), apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating event handler Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update event handler, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated event handler to update state
	// populated.
	apiResponse, err := r.client.GetEventHandler(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading event handler Info",
			"Could not read event handler Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	// Update resource state with updated items
	configDataModel := eventHandlerResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert event handler, unexpected error: "+err.Error(),
		)
		return
	}

	configDataModel.EventEndpointType = types.StringValue(apiResponse.EventEndpointType.Moniker)
	plan = configDataModel

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *eventHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state eventHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing event handler
	result, err := r.client.DeleteEventHandler(state.Moniker.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting event handler",
			"Could not delete event handler, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting event handler",
			"Could not delete event handler, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
