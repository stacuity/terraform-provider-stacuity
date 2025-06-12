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
	_ resource.Resource                = &endpointGroupResource{}
	_ resource.ResourceWithConfigure   = &endpointGroupResource{}
	_ resource.ResourceWithImportState = &endpointGroupResource{}
)

// NewEndpointGroupResource is a helper function to simplify the provider implementation.
func NewEndpointGroupResource() resource.Resource {
	return &endpointGroupResource{}
}

// endpointGroupResource is the resource implementation.
type endpointGroupResource struct {
	client *stacuity.Client
}

type endpointGroupResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Moniker               types.String `tfsdk:"moniker"`
	VSlice                types.String `tfsdk:"vslice"`
	EventMap              types.String `tfsdk:"event_map"`
	RoutingPolicy         types.String `tfsdk:"routing_policy"`
	SteeringProfile       types.String `tfsdk:"operator_policy"`
	RegionalGatewayPolicy types.String `tfsdk:"regional_gateway_policy"`
	IPAllocationType      types.String `tfsdk:"ip_allocation_type"`
}

func (r *endpointGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *endpointGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint_group"
}

func (r *endpointGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "endpointGroup resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the Endpoint Group.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the Endpoint Group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the Endpoint Group.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"vslice": schema.StringAttribute{
				Description: "The VSlice moniker that the Endpoint Group should use.",
				Required:    true,
			},
			"event_map": schema.StringAttribute{
				Description: "The Event Map API Moniker of that the Endpoint Group should use",
				Optional:    true,
			},
			"routing_policy": schema.StringAttribute{
				Description: "The Routing Policy Moniker that the Endpoint Group should use",
				Optional:    true,
			},
			"operator_policy": schema.StringAttribute{
				Description: "The operator policy Moniker that the Endpoint Group should use",
				Optional:    true,
			},
			"regional_gateway_policy": schema.StringAttribute{
				Description: "The Regional Gateway Policy Moniker that the Endpoint Group should use",
				Required:    true,
			},
			"ip_allocation_type": schema.StringAttribute{
				Description: "The IP allocation type that the Endpoint Group should use, such as static.",
				Required:    true,
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *endpointGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *endpointGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan endpointGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiData := models.EndpointGroupModifyItem{}
	var err = stacuity.ConvertToAPI(plan, &apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint group",
			"Could not create API endpoint group, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new endpointGroup
	createResponse, err := r.client.CreateEndpointGroup(apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating endpoint group",
			"Could not create endpoint group, unexpected error: "+err.Error(),
		)
		return
	}

	if createResponse.Success {
		getResponse, err := r.client.GetEndpointGroup(createResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading endpoint group",
				"Could not read endpoint group, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel := endpointGroupResourceModel{}
		err = stacuity.ConvertFromAPI(getResponse, &configDataModel)
		configDataModel.VSlice = types.StringValue(getResponse.VSlice.Moniker)
		configDataModel.IPAllocationType = types.StringValue(getResponse.IPAllocationType.Moniker)
		configDataModel.RegionalGatewayPolicy = types.StringValue(getResponse.RegionalGatewayPolicy.Moniker)

		if getResponse.EventMap != nil {
			configDataModel.EventMap = types.StringValue(getResponse.EventMap.Moniker)
		}
		if getResponse.RoutingPolicy != nil {
			configDataModel.RoutingPolicy = types.StringValue(getResponse.RoutingPolicy.Moniker)
		}
		if getResponse.SteeringProfile != nil {
			configDataModel.SteeringProfile = types.StringValue(getResponse.SteeringProfile.Moniker)
		}

		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting to TF",
				"Could not convert from API endpoint group, unexpected error: "+err.Error(),
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
func (r *endpointGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state endpointGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed endpointGroup values
	apiResponse, err := r.client.GetEndpointGroup(state.Moniker.ValueString())

	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading endpoint group Info",
			"Could not read endpoint group Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	configDataModel := endpointGroupResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert endpoint group, unexpected error: "+err.Error(),
		)
		return
	}
	configDataModel.VSlice = types.StringValue(apiResponse.VSlice.Moniker)
	configDataModel.IPAllocationType = types.StringValue(apiResponse.IPAllocationType.Moniker)
	configDataModel.RegionalGatewayPolicy = types.StringValue(apiResponse.RegionalGatewayPolicy.Moniker)
	if apiResponse.EventMap != nil {
		configDataModel.EventMap = types.StringValue(apiResponse.EventMap.Moniker)
	}
	if apiResponse.RoutingPolicy != nil {
		configDataModel.RoutingPolicy = types.StringValue(apiResponse.RoutingPolicy.Moniker)
	}
	if apiResponse.SteeringProfile != nil {
		configDataModel.SteeringProfile = types.StringValue(apiResponse.SteeringProfile.Moniker)
	}
	state = configDataModel

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *endpointGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan endpointGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiConfigDataModel := models.EndpointGroupModifyItem{}
	err := stacuity.ConvertToAPI(plan, &apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert endpoint group, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing endpointGroup
	_, err = r.client.UpdateEndpointGroup(plan.Moniker.ValueString(), apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating endpoint group Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update endpoint group, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated endpointGroup to update state
	// populated.
	apiResponse, err := r.client.GetEndpointGroup(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading endpoint group Info",
			"Could not read endpoint group Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	// Update resource state with updated items
	configDataModel := endpointGroupResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert endpoint group, unexpected error: "+err.Error(),
		)
		return
	}
	configDataModel.VSlice = types.StringValue(apiResponse.VSlice.Moniker)
	configDataModel.IPAllocationType = types.StringValue(apiResponse.IPAllocationType.Moniker)
	configDataModel.RegionalGatewayPolicy = types.StringValue(apiResponse.RegionalGatewayPolicy.Moniker)

	if apiResponse.EventMap != nil {
		configDataModel.EventMap = types.StringValue(apiResponse.EventMap.Moniker)
	}
	if apiResponse.RoutingPolicy != nil {
		configDataModel.RoutingPolicy = types.StringValue(apiResponse.RoutingPolicy.Moniker)
	}
	if apiResponse.SteeringProfile != nil {
		configDataModel.SteeringProfile = types.StringValue(apiResponse.SteeringProfile.Moniker)
	}
	plan = configDataModel

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *endpointGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state endpointGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing endpointGroup
	result, err := r.client.DeleteEndpointGroup(state.Moniker.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting endpoint group",
			"Could not delete endpoint group, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting endpoint group",
			"Could not delete endpoint group, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
