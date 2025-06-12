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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	stacuity "stacuity.com/go_client"
	models "stacuity.com/go_client/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vSliceResource{}
	_ resource.ResourceWithConfigure   = &vSliceResource{}
	_ resource.ResourceWithImportState = &vSliceResource{}
)

// NewVSliceResource is a helper function to simplify the provider implementation.
func NewVSliceResource() resource.Resource {
	return &vSliceResource{}
}

// vSliceResource is the resource implementation.
type vSliceResource struct {
	client *stacuity.Client
}

type vSlicesResourceModel struct {
	Id               types.String   `tfsdk:"id"`
	Name             types.String   `tfsdk:"name"`
	Moniker          types.String   `tfsdk:"moniker"`
	DNSServers       []types.String `tfsdk:"dns_servers"`
	DNSMode          types.String   `tfsdk:"dns_mode"`
	IpAddressFamily  types.String   `tfsdk:"ip_address_family"`
	SubnetAddress    types.String   `tfsdk:"subnet_address"`
	IpAllocationType types.String   `tfsdk:"ip_allocation_type"`
	EventMap         types.String   `tfsdk:"event_map"`
}

func (r *vSliceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vSliceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vslice"
}

func (r vSliceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data vSlicesResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If DNSMode is not configured, return without warning.
	if data.DNSMode.IsNull() {
		return
	}

	var dnsMode = data.DNSMode.ValueString()
	if dnsMode != "custom" && dnsMode != "auto" {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_mode"),
			"Invalid Configuration",
			"dns_mode must be either 'static' or 'auto'",
		)
	}

	var dnsServerCount = len(data.DNSServers)
	if dnsMode == "auto" && dnsServerCount == 0 {
		return
	} else if dnsMode == "auto" && dnsServerCount > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_servers"),
			"Invalid Configuration",
			"dns_servers shoudn't have a value if using dns_mode auto",
		)
	}

	if dnsServerCount == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("dns_servers"),
			"Invalid Configuration",
			"dns_servers must be specified with 'custom' dns_mode set",
		)
	}
}

func (r *vSliceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "vSlice resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the vSlice.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the vSlice",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the vSlice",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"subnet_address": schema.StringAttribute{
				Description: "Subnet applied to vSlice. This is the initial subnet, you can add more subnets after the VSlice has been created.",
				Required:    true,
			},
			"ip_allocation_type": schema.StringAttribute{
				Description: "Type of ip allocated Static/Pooled.",
				Default:     stringdefault.StaticString("static"),
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_mode": schema.StringAttribute{
				Description: "Type of DNS. Auto or Custom",
				Required:    true,
			},
			"dns_servers": schema.ListAttribute{
				Description: "DNS servers applied to vSlice, if using custom DNS",
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
			},
			"event_map": schema.StringAttribute{
				Description: "The moniker of the eventmap to link",
				Required:    false,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Computed:    true,
			},
			"ip_address_family": schema.StringAttribute{
				Description: "The type of IP address. Ipv4 or Ipv6",
				Optional:    true,
				Default:     stringdefault.StaticString("Ipv4"),
				Computed:    true,
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *vSliceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vSliceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vSlicesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	vSlice := models.VSliceModifyItem{
		Name:             plan.Name.ValueString(),
		Moniker:          plan.Moniker.ValueString(),
		EventMap:         plan.EventMap.ValueString(),
		DNSMode:          plan.DNSMode.ValueString(),
		IpAddressFamily:  plan.IpAddressFamily.ValueString(),
		IpAllocationType: plan.IpAllocationType.ValueString(),
		SubnetAddress:    plan.SubnetAddress.ValueString(),
	}

	for _, dns := range plan.DNSServers {
		vSlice.DNSServers = append(vSlice.DNSServers, dns.ValueString())
	}

	// Create new vSlice
	apiResponse, err := r.client.CreateVSlice(vSlice)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vSlice",
			"Could not create vSlice, unexpected error: "+err.Error(),
		)
		return
	}

	if apiResponse.Success {
		getResponse, err := r.client.GetVSlice(apiResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading vSlice",
				"Could not read vSlice, unexpected error: "+err.Error(),
			)
			return
		}

		plan = vSlicesResourceModel{
			Id:               types.StringValue(getResponse.Id),
			Moniker:          types.StringValue(getResponse.Moniker),
			Name:             types.StringValue(getResponse.Name),
			EventMap:         types.StringValue(getResponse.EventMap.Moniker),
			DNSMode:          types.StringValue(getResponse.DNSMode.Moniker),
			IpAddressFamily:  types.StringValue(getResponse.IpAddressFamily.Moniker),
			IpAllocationType: types.StringValue(plan.IpAllocationType.ValueString()),
		}

		if len(getResponse.Subnets) > 0 {
			plan.SubnetAddress = types.StringValue(getResponse.Subnets[0])
		}

		for _, dns := range getResponse.DNSServers {
			plan.DNSServers = append(plan.DNSServers, types.StringValue(dns))
		}
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
func (r *vSliceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state vSlicesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed vSlice values
	apiResponse, err := r.client.GetVSlice(state.Moniker.ValueString())
	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading vSlice Info",
			"Could not read vSlice Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Id = types.StringValue(apiResponse.Id)
	state.Moniker = types.StringValue(apiResponse.Moniker)
	state.Name = types.StringValue(apiResponse.Name)
	state.EventMap = types.StringValue(apiResponse.EventMap.Moniker)
	state.DNSMode = types.StringValue(apiResponse.DNSMode.Moniker)
	state.IpAddressFamily = types.StringValue(apiResponse.IpAddressFamily.Moniker)

	if len(apiResponse.Subnets) > 0 {
		state.SubnetAddress = types.StringValue(apiResponse.Subnets[0])
	}

	state.DNSServers = nil
	if len(apiResponse.DNSServers) > 0 && len(apiResponse.DNSServers[0]) > 0 {
		for _, dns := range apiResponse.DNSServers {
			state.DNSServers = append(state.DNSServers, types.StringValue(dns))
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vSliceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan vSlicesResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	vSlice := models.VSliceModifyItem{
		Id:              plan.Id.ValueString(),
		Name:            plan.Name.ValueString(),
		Moniker:         plan.Moniker.ValueString(),
		EventMap:        plan.EventMap.ValueString(),
		DNSMode:         plan.DNSMode.ValueString(),
		IpAddressFamily: plan.IpAddressFamily.ValueString(),
		SubnetAddress:   plan.SubnetAddress.ValueString(),
	}

	for _, dns := range plan.DNSServers {
		vSlice.DNSServers = append(vSlice.DNSServers, dns.ValueString())
	}

	// Update existing vSlice
	_, err := r.client.UpdateVSlice(plan.Moniker.ValueString(), vSlice)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating vSlice Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update vSlice, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated vSlice to update state
	// populated.
	apiResponse, err := r.client.GetVSlice(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading vSlice Info",
			"Could not read vSlice Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	// Update resource state with updated items
	plan.Id = types.StringValue(apiResponse.Id)
	plan.Moniker = types.StringValue(apiResponse.Moniker)
	plan.Name = types.StringValue(apiResponse.Name)
	plan.EventMap = types.StringValue(apiResponse.EventMap.Moniker)
	plan.DNSMode = types.StringValue(apiResponse.DNSMode.Moniker)
	plan.IpAddressFamily = types.StringValue(apiResponse.IpAddressFamily.Moniker)

	if len(apiResponse.Subnets) > 0 {
		plan.SubnetAddress = types.StringValue(apiResponse.Subnets[0])
	}

	for _, dns := range apiResponse.DNSServers {
		plan.DNSServers = append(plan.DNSServers, types.StringValue(dns))
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vSliceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vSlicesResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing vSlice
	result, err := r.client.DeleteVSlice(state.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting vSlice",
			"Could not delete vSlice, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting vSlice",
			"Could not delete vSlice, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
