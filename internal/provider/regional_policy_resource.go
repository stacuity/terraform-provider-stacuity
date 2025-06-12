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
	_ resource.Resource                = &regionalPolicyResource{}
	_ resource.ResourceWithConfigure   = &regionalPolicyResource{}
	_ resource.ResourceWithImportState = &regionalPolicyResource{}
)

// NewRegionalPolicyResource is a helper function to simplify the provider implementation.
func NewRegionalPolicyResource() resource.Resource {
	return &regionalPolicyResource{}
}

// regionalPolicyResource is the resource implementation.
type regionalPolicyResource struct {
	client *stacuity.Client
}

type regionalPolicyResourceModel struct {
	Id      types.String                   `tfsdk:"id"`
	Moniker types.String                   `tfsdk:"moniker"`
	Name    types.String                   `tfsdk:"name"`
	Entries *[]regionalPolicyEntryResource `tfsdk:"entries"`
}

type regionalPolicyEntryResource struct {
	RegionalGatewayId types.String `tfsdk:"regional_gateway_id"`
	OperatorId        types.Int32  `tfsdk:"operator_id"`
	Iso3              types.String `tfsdk:"iso_3"`
}

func (r *regionalPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *regionalPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regional_policy"
}

func hasDefaultEntry(entries []regionalPolicyEntryResource) bool {
	for _, entry := range entries {
		if entry.Iso3.IsNull() && entry.OperatorId.IsNull() {
			return true
		}
	}
	return false
}

func (r regionalPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data regionalPolicyResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If no entries are configured, return without warning.
	if data.Entries == nil {
		return
	}

	if !hasDefaultEntry(*data.Entries) {
		resp.Diagnostics.AddAttributeError(
			path.Root("entries"),
			"Invalid Configuration",
			"You must specify a default regional policy gateway. Create a new entry without ISO_3 and OperatorId.'",
		)
	}
}

func (r *regionalPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "regional policy resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the Regional Policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the Regional Policy.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the Regional Policy.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"entries": schema.SetNestedAttribute{
				Description: "List of entry rules attached to regional policy.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"regional_gateway_id": schema.StringAttribute{
							Description: "The regional gateway id/moniker to apply to",
							Optional:    true,
						},
						"operator_id": schema.Int32Attribute{
							Description: "The operator id to apply to",
							Optional:    true,
						},
						"iso_3": schema.StringAttribute{
							Description: "The iso id the rule applies to",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *regionalPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *regionalPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan regionalPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiData := models.RegionalPolicyModifyItem{}
	var err = stacuity.ConvertToAPI(plan, &apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating regional policy",
			"Could not create API regional policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new regional policy
	createResponse, err := r.client.CreateRegionalPolicy(apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating regional policy",
			"Could not create regional policy, unexpected error: "+err.Error(),
		)
		return
	}

	if createResponse.Success {

		getResponse, err := r.client.GetRegionalPolicy(createResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading regional policy",
				"Could not read regional policy, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel := regionalPolicyResourceModel{}

		// Add Entries
		if apiData.Entries != nil && len(*apiData.Entries) > 0 {
			_, err := r.client.AddRegionalPolicyEntries(*apiData.Entries, getResponse.Moniker)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating regional policy",
					"Could not create regional policy, unexpected error: "+err.Error(),
				)
				return
			}

			getEntriesResponse, err := r.client.GetRegionalPolicyEntries(getResponse.Moniker)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting regional policy",
					"Could not getting regional policy, unexpected error: "+err.Error(),
				)
				return
			}

			if len(getEntriesResponse) > 0 {
				getResponse.Entries = &getEntriesResponse
			}
		}

		err = stacuity.ConvertFromAPI(getResponse, &configDataModel)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error converting to TF",
				"Could not convert from API regional policy, unexpected error: "+err.Error(),
			)
			return
		}

		//reset
		if getResponse.Entries != nil {
			configDataModel.Entries = &[]regionalPolicyEntryResource{}
			for _, entry := range *getResponse.Entries {
				newEntry := regionalPolicyEntryResource{}
				newEntry.Iso3 = types.StringPointerValue(entry.Iso3)
				newEntry.RegionalGatewayId = types.StringPointerValue(&entry.RegionalGateway.Moniker)
				newEntry.OperatorId = types.Int32PointerValue(entry.OperatorId)
				*configDataModel.Entries = append(*configDataModel.Entries, newEntry)
			}
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
func (r *regionalPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state regionalPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed regional policy values
	apiResponse, err := r.client.GetRegionalPolicy(state.Moniker.ValueString())

	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading regional policy Info",
			"Could not read regional policy Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	entryResponse, err := r.client.GetRegionalPolicyEntries(state.Moniker.ValueString())
	if err != nil {

		resp.Diagnostics.AddError(
			"Error Reading regional policy Entries",
			"Could not read regional policy Entries, unexpected error: "+err.Error(),
		)
		return
	}

	if len(entryResponse) > 0 {
		apiResponse.Entries = &entryResponse
	}

	configDataModel := regionalPolicyResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert regional policy, unexpected error: "+err.Error(),
		)
		return
	}

	//reset
	if apiResponse.Entries != nil {
		configDataModel.Entries = &[]regionalPolicyEntryResource{}
		for _, entry := range *apiResponse.Entries {
			newEntry := regionalPolicyEntryResource{}
			newEntry.Iso3 = types.StringPointerValue(entry.Iso3)
			newEntry.RegionalGatewayId = types.StringPointerValue(&entry.RegionalGateway.Moniker)
			newEntry.OperatorId = types.Int32PointerValue(entry.OperatorId)
			*configDataModel.Entries = append(*configDataModel.Entries, newEntry)
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
func (r *regionalPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan regionalPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiConfigDataModel := models.RegionalPolicyModifyItem{}
	err := stacuity.ConvertToAPI(plan, &apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert regional policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing regional policy
	_, err = r.client.UpdateRegionalPolicy(plan.Moniker.ValueString(), apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating regional policy Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update regional policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Add entries
	if apiConfigDataModel.Entries != nil && len(*apiConfigDataModel.Entries) > 0 {
		_, err = r.client.AddRegionalPolicyEntries(*apiConfigDataModel.Entries, plan.Moniker.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error adding regional policy entries Info Moniker:"+plan.Moniker.ValueString(),
				"Could not update regional policy entries, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Fetch updated regional policy to update state
	// populated.
	apiResponse, err := r.client.GetRegionalPolicy(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading regional policy Info",
			"Could not read regional policy Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	entriesResponse, err := r.client.GetRegionalPolicyEntries(plan.Moniker.ValueString())
	if err != nil {

		resp.Diagnostics.AddError(
			"Error Reading regional policy Entries",
			"Could not read regional policy entries, unexpected error: "+err.Error(),
		)
		return
	}

	if len(entriesResponse) > 0 {
		apiResponse.Entries = &entriesResponse
	}

	// Update resource state with updated items
	configDataModel := regionalPolicyResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert regional policy, unexpected error: "+err.Error(),
		)
		return
	}

	//reset
	if apiResponse.Entries != nil {
		configDataModel.Entries = &[]regionalPolicyEntryResource{}
		for _, entry := range *apiResponse.Entries {
			newEntry := regionalPolicyEntryResource{}
			newEntry.Iso3 = types.StringPointerValue(entry.Iso3)
			newEntry.RegionalGatewayId = types.StringPointerValue(&entry.RegionalGateway.Moniker)
			newEntry.OperatorId = types.Int32PointerValue(entry.OperatorId)
			*configDataModel.Entries = append(*configDataModel.Entries, newEntry)
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
func (r *regionalPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state regionalPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing regional policy
	result, err := r.client.DeleteRegionalPolicy(state.Moniker.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting regional policy",
			"Could not delete regional policy, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting regional policy",
			"Could not delete regional policy, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
