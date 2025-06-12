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
	_ resource.Resource                = &routingPolicyResource{}
	_ resource.ResourceWithConfigure   = &routingPolicyResource{}
	_ resource.ResourceWithImportState = &routingPolicyResource{}
)

// NewRoutingPolicyResource is a helper function to simplify the provider implementation.
func NewRoutingPolicyResource() resource.Resource {
	return &routingPolicyResource{}
}

// routingPolicyResource is the resource implementation.
type routingPolicyResource struct {
	client *stacuity.Client
}

type routingPolicyResourceModel struct {
	Id                              types.String        `tfsdk:"id"`
	Name                            types.String        `tfsdk:"name"`
	Moniker                         types.String        `tfsdk:"moniker"`
	VSlice                          types.String        `tfsdk:"vslice"`
	RoutingPolicyStatus             types.String        `tfsdk:"routing_policy_status"`
	RoutingPolicyRules              []*RoutingRuleModel `tfsdk:"routing_policy_rules"`
	RoutingPolicyEdgeServices       []*EdgeServiceModel `tfsdk:"routing_policy_edge_services"`
	RateLimitUplinkMoniker          types.String        `tfsdk:"rate_limit_uplink_moniker"`
	RateLimitDownlinkMoniker        types.String        `tfsdk:"rate_limit_downlink_moniker"`
	PacketDiscardUplinkPercentage   types.Int32         `tfsdk:"packet_discard_uplink_percentage"`
	PacketDiscardDownlinkPercentage types.Int32         `tfsdk:"packet_discard_downlink_percentage"`
}

type RoutingRuleModel struct {
	Description            types.String `tfsdk:"description"`
	RuleAction             types.String `tfsdk:"rule_action"`
	RuleDirection          types.String `tfsdk:"rule_direction"`
	SourceIpPattern        types.String `tfsdk:"source_ip_pattern"`
	DestinationIpPattern   types.String `tfsdk:"destination_ip_pattern"`
	DivertIp               types.String `tfsdk:"divert_ip"`
	DivertPort             types.String `tfsdk:"divert_port"`
	TransportProtocol      types.String `tfsdk:"transport_protocol"`
	SourcePortPattern      types.String `tfsdk:"source_port_pattern"`
	DestinationPortPattern types.String `tfsdk:"destination_port_pattern"`
	RoutingTarget          types.String `tfsdk:"routing_target"`
	Reflexive              types.Bool   `tfsdk:"reflexive"`
	RegionalGateway        types.String `tfsdk:"regional_gateway"`
	Enabled                types.Bool   `tfsdk:"enabled"`
}

type EdgeServiceModel struct {
	Moniker                types.String   `tfsdk:"moniker"`
	Enabled                types.Bool     `tfsdk:"enabled"`
	EdgeServiceInstanceIds []types.String `tfsdk:"edge_service_instance_ids"`
}

func (r *routingPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *routingPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_policy"
}

func (r *routingPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "routing policy resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the routing policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the routing policy.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(150),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the routing policy.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"vslice": schema.StringAttribute{
				Description: "The VSlice that the routing policy belongs to.",
				Required:    true,
			},
			"routing_policy_status": schema.StringAttribute{
				Description: "Status of the routing policy.",
				Required:    true,
			},
			"routing_policy_rules": schema.SetNestedAttribute{
				Description: "List of rules for the routing policy.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Description: "Description of the rule.",
							Required:    true,
						},
						"rule_action": schema.StringAttribute{
							Description: "The action to take on packets that match this rule.",
							Required:    true,
						},
						"rule_direction": schema.StringAttribute{
							Description: "Direction of traffic for the rule.",
							Required:    true,
						},
						"source_ip_pattern": schema.StringAttribute{
							Description: "IP pattern for source IPs.",
							Optional:    true,
						},
						"destination_ip_pattern": schema.StringAttribute{
							Description: "IP pattern for destination IPs.",
							Optional:    true,
						},
						"divert_ip": schema.StringAttribute{
							Description: "IP address to divert traffic to.",
							Optional:    true,
						},
						"divert_port": schema.StringAttribute{
							Description: "Port number to divert traffic to.",
							Optional:    true,
						},
						"transport_protocol": schema.StringAttribute{
							Description: "Transport protocol for the rule.",
							Optional:    true,
						},
						"source_port_pattern": schema.StringAttribute{
							Description: "Port pattern for source ports.",
							Optional:    true,
						},
						"destination_port_pattern": schema.StringAttribute{
							Description: "Port pattern for destination ports.",
							Optional:    true,
						},
						"routing_target": schema.StringAttribute{
							Description: "The routing target for the rule.",
							Optional:    true,
						},
						"reflexive": schema.BoolAttribute{
							Description: "Whether this is a reflexive rule.",
							Required:    true,
						},
						"regional_gateway": schema.StringAttribute{
							Description: "Regional gateway for the rule.",
							Optional:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether this rule is enabled.",
							Optional:    true,
						},
					},
				},
			},
			"routing_policy_edge_services": schema.SetNestedAttribute{
				Description: "List of edge services for the routing policy.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"moniker": schema.StringAttribute{
							Description: "API Moniker of the edge service.",
							Optional:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether this edge service is enabled.",
							Optional:    true,
						},
						"edge_service_instance_ids": schema.ListAttribute{
							ElementType: types.StringType,
							Description: "API Monikers of the edge service instance.",
							Optional:    true,
						},
					},
				},
			},
			"rate_limit_uplink_moniker": schema.StringAttribute{
				Description: "API Moniker for the uplink rate limit.",
				Optional:    true,
			},
			"rate_limit_downlink_moniker": schema.StringAttribute{
				Description: "API Moniker for the downlink rate limit.",
				Optional:    true,
			},
			"packet_discard_uplink_percentage": schema.Int32Attribute{
				Description: "Percentage of uplink packets to discard.",
				Optional:    true,
			},
			"packet_discard_downlink_percentage": schema.Int32Attribute{
				Description: "Percentage of downlink packets to discard.",
				Optional:    true,
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *routingPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *routingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan routingPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiData := models.RoutingPolicyModifyItem{}
	var err = stacuity.ConvertToAPI(plan, &apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating routing policy",
			"Could not create API routing policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new routing policy
	apiResponse, err := r.client.CreateRoutingPolicy(apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating routing policy",
			"Could not create routing policy, unexpected error: "+err.Error(),
		)
		return
	}

	if apiResponse.Success {
		getResponse, err := r.client.GetRoutingPolicy(apiResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading routing policy",
				"Could not read routing policy, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel := routingPolicyResourceModel{}
		err = stacuity.ConvertFromAPI(getResponse, &configDataModel)

		configDataModel.RateLimitDownlinkMoniker = types.StringValue(getResponse.RateLimitDownlink.Moniker)
		configDataModel.RateLimitUplinkMoniker = types.StringValue(getResponse.RateLimitUplink.Moniker)
		configDataModel.RoutingPolicyStatus = types.StringValue(getResponse.RoutingPolicyStatus.Moniker)
		configDataModel.VSlice = types.StringValue(getResponse.VSlice.Moniker)

		for id, rule := range *getResponse.RoutingPolicyRules {
			mappedRule := configDataModel.RoutingPolicyRules[id]
			mappedRule.RuleDirection = types.StringValue(rule.RuleDirection.Moniker)
			mappedRule.RuleAction = types.StringValue(rule.RuleAction.Moniker)
			if rule.RegionalGateway != nil {
				mappedRule.RegionalGateway = types.StringValue(rule.RegionalGateway.Moniker)
			}

			if rule.TransportProtocol != nil {
				mappedRule.TransportProtocol = types.StringValue(rule.TransportProtocol.Moniker)
			}

			if rule.RoutingTarget != nil {
				mappedRule.RoutingTarget = types.StringValue(rule.RoutingTarget.Moniker)
			}

			if rule.DivertPort != nil && *rule.DivertPort == "" {
				mappedRule.DivertPort = types.StringNull()
			}

			configDataModel.RoutingPolicyRules[id] = mappedRule
		}

		//reset
		configDataModel.RoutingPolicyEdgeServices = nil
		for _, edge := range *getResponse.RoutingPolicyEdgeServices {
			for _, planEdge := range plan.RoutingPolicyEdgeServices {
				if edge.Moniker == planEdge.Moniker.ValueString() {
					planEdge.Enabled = types.BoolValue(edge.Enabled)
					configDataModel.RoutingPolicyEdgeServices = append(configDataModel.RoutingPolicyEdgeServices, planEdge)
				}
			}
		}

		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting to TF",
				"Could not convert from API routing policy, unexpected error: "+err.Error(),
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
func (r *routingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state routingPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed routing policy values
	apiResponse, err := r.client.GetRoutingPolicy(state.Moniker.ValueString())

	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading routing policy Info",
			"Could not read routing policy Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	configDataModel := routingPolicyResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert routing policy, unexpected error: "+err.Error(),
		)
		return
	}

	configDataModel.RateLimitDownlinkMoniker = types.StringValue(apiResponse.RateLimitDownlink.Moniker)
	configDataModel.RateLimitUplinkMoniker = types.StringValue(apiResponse.RateLimitUplink.Moniker)
	configDataModel.RoutingPolicyStatus = types.StringValue(apiResponse.RoutingPolicyStatus.Moniker)
	configDataModel.VSlice = types.StringValue(apiResponse.VSlice.Moniker)

	for id, rule := range *apiResponse.RoutingPolicyRules {
		mappedRule := configDataModel.RoutingPolicyRules[id]
		mappedRule.RuleDirection = types.StringValue(rule.RuleDirection.Moniker)
		mappedRule.RuleAction = types.StringValue(rule.RuleAction.Moniker)

		if rule.RegionalGateway != nil {
			mappedRule.RegionalGateway = types.StringValue(rule.RegionalGateway.Moniker)
		}

		if rule.TransportProtocol != nil {
			mappedRule.TransportProtocol = types.StringValue(rule.TransportProtocol.Moniker)
		}

		if rule.RoutingTarget != nil {
			mappedRule.RoutingTarget = types.StringValue(rule.RoutingTarget.Moniker)
		}

		if rule.DivertPort != nil && *rule.DivertPort == "" {
			mappedRule.DivertPort = types.StringNull()
		}

		configDataModel.RoutingPolicyRules[id] = mappedRule
	}

	//reset
	configDataModel.RoutingPolicyEdgeServices = nil
	for _, edge := range *apiResponse.RoutingPolicyEdgeServices {
		for _, planEdge := range state.RoutingPolicyEdgeServices {
			if edge.Moniker == planEdge.Moniker.ValueString() {
				planEdge.Enabled = types.BoolValue(edge.Enabled)
				configDataModel.RoutingPolicyEdgeServices = append(configDataModel.RoutingPolicyEdgeServices, planEdge)
			}
		}
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
func (r *routingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan routingPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiConfigDataModel := models.RoutingPolicyModifyItem{}
	err := stacuity.ConvertToAPI(plan, &apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert routing policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing routing policy
	_, err = r.client.UpdateRoutingPolicy(plan.Moniker.ValueString(), apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating routing policy Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update routing policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated routing policy to update state
	// populated.
	apiResponse, err := r.client.GetRoutingPolicy(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading routing policy Info",
			"Could not read routing policy Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	// Update resource state with updated items
	configDataModel := routingPolicyResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert routing policy, unexpected error: "+err.Error(),
		)
		return
	}

	configDataModel.RateLimitDownlinkMoniker = types.StringValue(apiResponse.RateLimitDownlink.Moniker)
	configDataModel.RateLimitUplinkMoniker = types.StringValue(apiResponse.RateLimitUplink.Moniker)
	configDataModel.RoutingPolicyStatus = types.StringValue(apiResponse.RoutingPolicyStatus.Moniker)
	configDataModel.VSlice = types.StringValue(apiResponse.VSlice.Moniker)

	for id, rule := range *apiResponse.RoutingPolicyRules {
		mappedRule := configDataModel.RoutingPolicyRules[id]
		mappedRule.RuleDirection = types.StringValue(rule.RuleDirection.Moniker)
		mappedRule.RuleAction = types.StringValue(rule.RuleAction.Moniker)
		if rule.RegionalGateway != nil {
			mappedRule.RegionalGateway = types.StringValue(rule.RegionalGateway.Moniker)
		}

		if rule.TransportProtocol != nil {
			mappedRule.TransportProtocol = types.StringValue(rule.TransportProtocol.Moniker)
		}

		if rule.RoutingTarget != nil {
			mappedRule.RoutingTarget = types.StringValue(rule.RoutingTarget.Moniker)
		}

		if rule.DivertPort != nil && *rule.DivertPort == "" {
			mappedRule.DivertPort = types.StringNull()
		}

		configDataModel.RoutingPolicyRules[id] = mappedRule
	}

	//reset
	configDataModel.RoutingPolicyEdgeServices = nil
	for _, edge := range *apiResponse.RoutingPolicyEdgeServices {
		for _, planEdge := range plan.RoutingPolicyEdgeServices {
			if edge.Moniker == planEdge.Moniker.ValueString() {
				planEdge.Enabled = types.BoolValue(edge.Enabled)
				configDataModel.RoutingPolicyEdgeServices = append(configDataModel.RoutingPolicyEdgeServices, planEdge)
			}
		}
	}

	plan = configDataModel

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *routingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state routingPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing routing policy
	result, err := r.client.DeleteRoutingPolicy(state.Moniker.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting routing policy",
			"Could not delete routing policy, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting routing policy",
			"Could not delete routing policy, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
