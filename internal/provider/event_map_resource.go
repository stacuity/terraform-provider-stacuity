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
	_ resource.Resource                = &eventMapResource{}
	_ resource.ResourceWithConfigure   = &eventMapResource{}
	_ resource.ResourceWithImportState = &eventMapResource{}
)

// NewEventMapResource is a helper function to simplify the provider implementation.
func NewEventMapResource() resource.Resource {
	return &eventMapResource{}
}

// eventMapResource is the resource implementation.
type eventMapResource struct {
	client *stacuity.Client
}

type eventMapResourceModel struct {
	Id            types.String            `tfsdk:"id"`
	Name          types.String            `tfsdk:"name"`
	Moniker       types.String            `tfsdk:"moniker"`
	EventScope    types.String            `tfsdk:"event_scope"`
	Subscriptions *[]subscriptionResource `tfsdk:"subscriptions"`
}

type subscriptionResource struct {
	EventEndpointId types.String `tfsdk:"event_endpoint_id"`
	EventTypeId     types.String `tfsdk:"event_type_id"`
}

type subscriptionResourceRead struct {
	EventEndpointId types.String `tfsdk:"event_endpoint_id"`
	EventTypeId     types.String `tfsdk:"event_type_id"`
	EventType       eventType    `tfsdk:"event_type"`
	EventHandler    eventHandler `tfsdk:"event_handler"`
}

type eventHandler struct {
	Moniker types.String `tfsdk:"moniker"`
	Name    types.String `tfsdk:"name"`
	Active  types.Bool   `tfsdk:"active"`
}

func (r *eventMapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *eventMapResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_map"
}

func (r *eventMapResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "event map resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the event map.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the event map.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the event map.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"event_scope": schema.StringAttribute{
				Description: "API Moniker of the event scope.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subscriptions": schema.SetNestedAttribute{
				Description: "List of subscriptions attached to event map.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"event_endpoint_id": schema.StringAttribute{
							Description: "The monkier of the event handler.",
							Required:    true,
						},
						"event_type_id": schema.StringAttribute{
							Description: "The type of event to subscribe to",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *eventMapResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *eventMapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan eventMapResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiData := models.EventMapModifyItem{}
	var err = stacuity.ConvertToAPI(plan, &apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating event map",
			"Could not create API event map, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new event map
	createResponse, err := r.client.CreateEventMap(apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating event map",
			"Could not create event map, unexpected error: "+err.Error(),
		)
		return
	}

	if createResponse.Success {

		getResponse, err := r.client.GetEventMap(createResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading event map",
				"Could not read event map, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel := eventMapResourceModel{}

		err = stacuity.ConvertFromAPI(getResponse, &configDataModel)

		// Add subscriptions
		if apiData.Subscriptions != nil && len(*apiData.Subscriptions) > 0 {
			createSubResponse, err := r.client.AddEventMapSubscriptions(*apiData.Subscriptions, getResponse.Moniker)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating event map",
					"Could not create event map, unexpected error: "+err.Error(),
				)
				return
			}

			subscriptionResourceRead := &[]subscriptionResourceRead{}
			err = stacuity.ConvertFromAPI(createSubResponse.Data, subscriptionResourceRead)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error converting to TF",
					"Could not convert event map, unexpected error: "+err.Error(),
				)
				return
			}

			configDataModel.Subscriptions = &[]subscriptionResource{}
			for _, sub := range *subscriptionResourceRead {
				newSub := subscriptionResource{}
				newSub.EventEndpointId = types.StringValue(sub.EventHandler.Moniker.ValueString())
				newSub.EventTypeId = types.StringValue(sub.EventType.Moniker.ValueString())
				*configDataModel.Subscriptions = append(*configDataModel.Subscriptions, newSub)
			}
		}

		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting to TF",
				"Could not convert from API event map subscription, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel.EventScope = types.StringValue(getResponse.EventScope.Moniker)
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
func (r *eventMapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state eventMapResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed event map values
	apiResponse, err := r.client.GetEventMap(state.Moniker.ValueString())

	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading event map Info",
			"Could not read event map Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	subscriptionsResponse, err := r.client.GetEventMapSubscriptions(state.Moniker.ValueString())
	if err != nil {

		resp.Diagnostics.AddError(
			"Error Reading event map Subscriptions",
			"Could not read event map subscriptions, unexpected error: "+err.Error(),
		)
		return
	}

	if len(subscriptionsResponse) > 0 {
		apiResponse.Subscriptions = &subscriptionsResponse
	}

	configDataModel := eventMapResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert event map, unexpected error: "+err.Error(),
		)
		return
	}

	//reset
	if apiResponse.Subscriptions != nil {
		configDataModel.Subscriptions = &[]subscriptionResource{}
		for _, sub := range *apiResponse.Subscriptions {
			newSub := subscriptionResource{}
			newSub.EventEndpointId = types.StringValue(sub.EventEndpoint.Moniker)
			newSub.EventTypeId = types.StringValue(sub.EventType.Moniker)
			*configDataModel.Subscriptions = append(*configDataModel.Subscriptions, newSub)
		}
	}

	configDataModel.EventScope = types.StringValue(apiResponse.EventScope.Moniker)
	state = configDataModel

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *eventMapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan eventMapResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiConfigDataModel := models.EventMapModifyItem{}
	err := stacuity.ConvertToAPI(plan, &apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert event map, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing event map
	_, err = r.client.UpdateEventMap(plan.Moniker.ValueString(), apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating event map Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update event map, unexpected error: "+err.Error(),
		)
		return
	}

	// Add subscriptions
	if apiConfigDataModel.Subscriptions != nil && len(*apiConfigDataModel.Subscriptions) > 0 {
		_, err = r.client.AddEventMapSubscriptions(*apiConfigDataModel.Subscriptions, plan.Moniker.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error adding event map subscriptions Info Moniker:"+plan.Moniker.ValueString(),
				"Could not update event map subscriptions, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Fetch updated event map to update state
	// populated.
	apiResponse, err := r.client.GetEventMap(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading event map Info",
			"Could not read event map Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	subscriptionsResponse, err := r.client.GetEventMapSubscriptions(plan.Moniker.ValueString())
	if err != nil {

		resp.Diagnostics.AddError(
			"Error Reading event map Subscriptions",
			"Could not read event map subscriptions, unexpected error: "+err.Error(),
		)
		return
	}

	subscriptionResourceRead := &[]subscriptionResourceRead{}
	if len(subscriptionsResponse) > 0 {
		err = stacuity.ConvertFromAPI(subscriptionsResponse, subscriptionResourceRead)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting to TF",
				"Could not convert event map, unexpected error: "+err.Error(),
			)
			return
		}

	}

	// Update resource state with updated items
	configDataModel := eventMapResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert event map, unexpected error: "+err.Error(),
		)
		return
	}

	//reset
	if len(*subscriptionResourceRead) > 0 {
		configDataModel.Subscriptions = &[]subscriptionResource{}
		for _, sub := range *subscriptionResourceRead {
			newSub := subscriptionResource{}
			newSub.EventEndpointId = types.StringValue(sub.EventHandler.Moniker.ValueString())
			newSub.EventTypeId = types.StringValue(sub.EventType.Moniker.ValueString())
			*configDataModel.Subscriptions = append(*configDataModel.Subscriptions, newSub)
		}
	}

	configDataModel.EventScope = types.StringValue(apiResponse.EventScope.Moniker)
	plan = configDataModel

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *eventMapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state eventMapResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing event map
	result, err := r.client.DeleteEventMap(state.Moniker.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting event map",
			"Could not delete event map, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting event map",
			"Could not delete event map, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
