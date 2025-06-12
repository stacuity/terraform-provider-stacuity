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
	_ resource.Resource                = &routingTargetResource{}
	_ resource.ResourceWithConfigure   = &routingTargetResource{}
	_ resource.ResourceWithImportState = &routingTargetResource{}
)

// NewRoutingTargetResource is a helper function to simplify the provider implementation.
func NewRoutingTargetResource() resource.Resource {
	return &routingTargetResource{}
}

// routingTargetResource is the resource implementation.
type routingTargetResource struct {
	client *stacuity.Client
}

type routingTargetResourceModel struct {
	Id                           types.String            `tfsdk:"id"`
	Name                         types.String            `tfsdk:"name"`
	Moniker                      types.String            `tfsdk:"moniker"`
	RoutingTargetType            types.String            `tfsdk:"routing_target_type"`
	RoutingRedundancyZoneMoniker types.String            `tfsdk:"redundancy_zone_moniker"`
	ConfigurationData            *ConfigurationDataModel `tfsdk:"configuration_data"`
	VSlice                       types.String            `tfsdk:"vslice"`
	RoutingTargetTypeInstanceId  types.String            `tfsdk:"routing_target_type_instance_id"`
}

type ConfigurationDataModel struct {
	VpnConfig       *RoutingTargetVpn       `tfsdk:"vpn_config"`
	WireGuardConfig *RoutingTargetWireguard `tfsdk:"wireguard_config"`
}

type RoutingTargetVpn struct {
	RemotePeerAddress      types.String ` tfsdk:"remote_peer_address"`
	RemoteSubnets          types.String `tfsdk:"remote_subnets"`
	RemoteEncryptionDomain types.String `tfsdk:"remote_encryption_domain"`
	LocalEncryptionDomain  types.String `tfsdk:"local_encryption_domain"`
	LocalSubnets           types.String `tfsdk:"local_subnets"`
	PresharedKey           types.String `tfsdk:"preshared_key"`
	KeyExchangeType        types.String `tfsdk:"key_exchange_type"`
	VpnIkeOption           types.String `tfsdk:"vpn_ike_option"`
	VpnEspOption           types.String `tfsdk:"vpn_esp_option"`
	Phase1Lifetime         types.Int32  `tfsdk:"phase1_lifetime"`
	Phase2Lifetime         types.Int32  `tfsdk:"phase2_lifetime"`
}

type RoutingTargetWireguard struct {
	LocalPublicKey       types.String `tfsdk:"local_public_key"`
	LocalSubnets         types.String `tfsdk:"local_subnets"`
	RemotePublicKey      types.String `tfsdk:"remote_public_key"`
	RemoteSubnets        types.String `tfsdk:"remote_subnets"`
	RemotePeerIPAddress  types.String `tfsdk:"remote_peer_ip_address"`
	RemotePeerPortNumber types.Int32  `tfsdk:"remote_peer_port_number"`
}

func (r *routingTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *routingTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_target"
}

func (r *routingTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "routing target resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the routing target.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the routing target",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the routing target",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"routing_target_type": schema.StringAttribute{
				Description: "The target type such as Internet, WireGuard or VPN",
				Required:    true,
			},
			"redundancy_zone_moniker": schema.StringAttribute{
				Description: "The moniker of the redundancy zone to use for the routing target",
				Required:    true,
			},
			"configuration_data": schema.SingleNestedAttribute{
				Description: "The configuration data for VPN or WireGuard settings.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"vpn_config": schema.SingleNestedAttribute{
						Description: "Configuration for RoutingTargetVpn",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"remote_peer_address": schema.StringAttribute{
								Description: "Remote peer address for VPN.",
								Optional:    true,
							},
							"remote_subnets": schema.StringAttribute{
								Description: "Remote subnets for VPN.",
								Optional:    true,
							},
							"remote_encryption_domain": schema.StringAttribute{
								Description: "Remote encryption domain for VPN.",
								Optional:    true,
							},
							"local_encryption_domain": schema.StringAttribute{
								Description: "Local encryption domain for VPN.",
								Optional:    true,
							},
							"local_subnets": schema.StringAttribute{
								Description: "Local subnets for VPN.",
								Optional:    true,
							},
							"preshared_key": schema.StringAttribute{
								Description: "Preshared key for VPN.",
								Optional:    true,
							},
							"key_exchange_type": schema.StringAttribute{
								Description: "Key exchange type for VPN.",
								Optional:    true,
							},
							"vpn_ike_option": schema.StringAttribute{
								Description: "VPN IKE option for VPN.",
								Optional:    true,
							},
							"vpn_esp_option": schema.StringAttribute{
								Description: "VPN ESP option for VPN.",
								Optional:    true,
							},
							"phase1_lifetime": schema.Int32Attribute{
								Description: "Phase 1 lifetime for VPN IKE negotiation.",
								Optional:    true,
							},
							"phase2_lifetime": schema.Int32Attribute{
								Description: "Phase 2 lifetime for VPN IKE negotiation.",
								Optional:    true,
							},
						},
					},
					"wireguard_config": schema.SingleNestedAttribute{
						Description: "Configuration for RoutingTargetWireguard",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"local_public_key": schema.StringAttribute{
								Description: "Local public key for WireGuard.",
								Optional:    true,
								Computed:    true,
							},
							"local_subnets": schema.StringAttribute{
								Description: "Local subnets for WireGuard.",
								Optional:    true,
							},
							"remote_public_key": schema.StringAttribute{
								Description: "Remote public key for WireGuard.",
								Optional:    true,
							},
							"remote_subnets": schema.StringAttribute{
								Description: "Remote subnets for WireGuard.",
								Optional:    true,
							},
							"remote_peer_ip_address": schema.StringAttribute{
								Description: "Remote peer IP address for WireGuard.",
								Optional:    true,
							},
							"remote_peer_port_number": schema.Int32Attribute{
								Description: "Remote peer port number for WireGuard.",
								Optional:    true,
								Computed:    true,
							},
						},
					},
				},
			},
			"vslice": schema.StringAttribute{
				Description: "The VSlice that the RoutingTarget belongs to.",
				Required:    true,
			},
			"routing_target_type_instance_id": schema.StringAttribute{
				Description: "Id or moniker of the routing target type instance",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *routingTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *routingTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan routingTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiData := models.RoutingTargetModifyItem{}
	var err = stacuity.ConvertToAPI(plan, &apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating routing target",
			"Could not create routing target, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new routing target
	apiResponse, err := r.client.CreateRoutingTarget(apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating routing target",
			"Could not create routing target, unexpected error: "+err.Error(),
		)
		return
	}

	if apiResponse.Success {
		getResponse, err := r.client.GetRoutingTarget(apiResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading routing target",
				"Could not read routing target, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel := routingTargetResourceModel{}
		err = stacuity.ConvertFromAPI(getResponse, &configDataModel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting to TF",
				"Could not convert routing target, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel.VSlice = types.StringValue(getResponse.VSlice.Moniker)
		configDataModel.RoutingTargetType = types.StringValue(getResponse.RoutingTargetType.Moniker)
		configDataModel.RoutingTargetTypeInstanceId = types.StringValue(getResponse.RoutingTargetTypeInstance.Moniker)
		handleEmptyConfigData(&configDataModel)
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
func (r *routingTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state routingTargetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed routing target values
	apiResponse, err := r.client.GetRoutingTarget(state.Moniker.ValueString())

	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading routing target Info",
			"Could not read routing target Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	configDataModel := routingTargetResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert routing target, unexpected error: "+err.Error(),
		)
		return
	}

	configDataModel.VSlice = types.StringValue(apiResponse.VSlice.Moniker)
	configDataModel.RoutingTargetType = types.StringValue(apiResponse.RoutingTargetType.Moniker)
	configDataModel.RoutingTargetTypeInstanceId = types.StringValue(apiResponse.RoutingTargetTypeInstance.Moniker)

	handleEmptyConfigData(&configDataModel)

	state = configDataModel

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func handleEmptyConfigData(configDataModel *routingTargetResourceModel) {
	if configDataModel.ConfigurationData.VpnConfig == nil && configDataModel.ConfigurationData.WireGuardConfig == nil {
		configDataModel.ConfigurationData = nil
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *routingTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan routingTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiConfigDataModel := models.RoutingTargetModifyItem{}
	err := stacuity.ConvertToAPI(plan, &apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert routing target, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing routing target
	_, err = r.client.UpdateRoutingTarget(plan.Moniker.ValueString(), apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating routing target Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update routing target, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated routing target to update state
	// populated.
	apiResponse, err := r.client.GetRoutingTarget(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading routing target Info",
			"Could not read routing target Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	// Update resource state with updated items
	configDataModel := routingTargetResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert routing target, unexpected error: "+err.Error(),
		)
		return
	}

	configDataModel.VSlice = types.StringValue(apiResponse.VSlice.Moniker)
	configDataModel.RoutingTargetType = types.StringValue(apiResponse.RoutingTargetType.Moniker)
	configDataModel.RoutingTargetTypeInstanceId = types.StringValue(apiResponse.RoutingTargetTypeInstance.Moniker)
	handleEmptyConfigData(&configDataModel)

	plan = configDataModel

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *routingTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state routingTargetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing routing target
	result, err := r.client.DeleteRoutingTarget(state.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting routing target",
			"Could not delete routing target, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting routing target",
			"Could not delete routing target, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
