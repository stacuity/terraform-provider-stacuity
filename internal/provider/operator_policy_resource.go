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
	_ resource.Resource                = &operatorPolicyResource{}
	_ resource.ResourceWithConfigure   = &operatorPolicyResource{}
	_ resource.ResourceWithImportState = &operatorPolicyResource{}
)

// NewOperatorPolicyResource is a helper function to simplify the provider implementation.
func NewOperatorPolicyResource() resource.Resource {
	return &operatorPolicyResource{}
}

// operatorPolicyResource is the resource implementation.
type operatorPolicyResource struct {
	client *stacuity.Client
}

type operatorPolicyResourceModel struct {
	Id      types.String                   `tfsdk:"id"`
	Moniker types.String                   `tfsdk:"moniker"`
	Name    types.String                   `tfsdk:"name"`
	Entries *[]operatorPolicyEntryResource `tfsdk:"entries"`
}

type operatorPolicyEntryResource struct {
	OperatorId                 types.Int32  `tfsdk:"operator_id"`
	Iso3                       types.String `tfsdk:"iso_3"`
	SteeringProfileEntryAction types.String `tfsdk:"steering_profile_entry_action"`
}

func (r *operatorPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import Id and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *operatorPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_operator_policy"
}

func (r *operatorPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "operator policy resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The identifier for the Operator Policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the Operator Policy.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"moniker": schema.StringAttribute{
				Description: "API Moniker of the Operator Policy.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(30),
				},
			},
			"entries": schema.SetNestedAttribute{
				Description: "List of entry rules attached to operator policy.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"operator_id": schema.Int32Attribute{
							Description: "The operator id to apply to",
							Optional:    true,
						},
						"iso_3": schema.StringAttribute{
							Description: "The iso id the rule applies to",
							Optional:    true,
						},
						"steering_profile_entry_action": schema.StringAttribute{
							Description: "The action to take for this network or country. allow, reject or reject-soft",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Configure implements resource.ResourceWithConfigure.
func (r *operatorPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *operatorPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan operatorPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	apiData := models.OperatorPolicyModifyItem{}
	var err = stacuity.ConvertToAPI(plan, &apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating operator policy",
			"Could not create API operator policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Create new operator policy
	createResponse, err := r.client.CreateOperatorPolicy(apiData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating operator policy",
			"Could not create operator policy, unexpected error: "+err.Error(),
		)
		return
	}

	if createResponse.Success {

		getResponse, err := r.client.GetOperatorPolicy(createResponse.Data)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error re-reading operator policy",
				"Could not read operator policy, unexpected error: "+err.Error(),
			)
			return
		}

		configDataModel := operatorPolicyResourceModel{}

		// Add Entries
		if apiData.Entries != nil && len(*apiData.Entries) > 0 {
			_, err := r.client.AddOperatorPolicyEntries(*apiData.Entries, getResponse.Moniker)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating operator policy",
					"Could not create operator policy, unexpected error: "+err.Error(),
				)
				return
			}

			getEntriesResponse, err := r.client.GetOperatorPolicyEntries(getResponse.Moniker)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting operator policy",
					"Could not getting operator policy, unexpected error: "+err.Error(),
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
				"Could not convert from API operator policy, unexpected error: "+err.Error(),
			)
			return
		}

		//reset
		if getResponse.Entries != nil {
			configDataModel.Entries = &[]operatorPolicyEntryResource{}
			for _, entry := range *getResponse.Entries {
				newEntry := operatorPolicyEntryResource{}
				newEntry.Iso3 = types.StringPointerValue(entry.Iso3)
				newEntry.OperatorId = types.Int32PointerValue(entry.OperatorId)
				newEntry.SteeringProfileEntryAction = types.StringValue(entry.SteeringProfileEntryAction.Moniker)
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
func (r *operatorPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state operatorPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get refreshed operator policy values
	apiResponse, err := r.client.GetOperatorPolicy(state.Moniker.ValueString())

	if err != nil {

		if err.Error() == "Record not found" {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading operator policy Info",
			"Could not read operator policy Moniker "+state.Moniker.ValueString()+": "+err.Error(),
		)
		return
	}

	entryResponse, err := r.client.GetOperatorPolicyEntries(state.Moniker.ValueString())
	if err != nil {

		resp.Diagnostics.AddError(
			"Error Reading operator policy Entries",
			"Could not read operator policy Entries, unexpected error: "+err.Error(),
		)
		return
	}

	if len(entryResponse) > 0 {
		apiResponse.Entries = &entryResponse
	}

	configDataModel := operatorPolicyResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert operator policy, unexpected error: "+err.Error(),
		)
		return
	}

	//reset
	if apiResponse.Entries != nil {
		configDataModel.Entries = &[]operatorPolicyEntryResource{}
		for _, entry := range *apiResponse.Entries {
			newEntry := operatorPolicyEntryResource{}
			newEntry.Iso3 = types.StringPointerValue(entry.Iso3)
			newEntry.OperatorId = types.Int32PointerValue(entry.OperatorId)
			newEntry.SteeringProfileEntryAction = types.StringValue(entry.SteeringProfileEntryAction.Moniker)
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
func (r *operatorPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan operatorPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiConfigDataModel := models.OperatorPolicyModifyItem{}
	err := stacuity.ConvertToAPI(plan, &apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert operator policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing operator policy
	_, err = r.client.UpdateOperatorPolicy(plan.Moniker.ValueString(), apiConfigDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating operator policy Info Moniker:"+plan.Moniker.ValueString(),
			"Could not update operator policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Add entries
	if apiConfigDataModel.Entries != nil && len(*apiConfigDataModel.Entries) > 0 {
		_, err = r.client.AddOperatorPolicyEntries(*apiConfigDataModel.Entries, plan.Moniker.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error adding operator policy entries Info Moniker:"+plan.Moniker.ValueString(),
				"Could not update operator policy entries, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Fetch updated operator policy to update state
	// populated.
	apiResponse, err := r.client.GetOperatorPolicy(plan.Moniker.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading operator policy Info",
			"Could not read operator policy Moniker "+plan.Moniker.ValueString()+" "+err.Error(),
		)
		return
	}

	entriesResponse, err := r.client.GetOperatorPolicyEntries(plan.Moniker.ValueString())
	if err != nil {

		resp.Diagnostics.AddError(
			"Error Reading operator policy Entries",
			"Could not read operator policy entries, unexpected error: "+err.Error(),
		)
		return
	}

	if len(entriesResponse) > 0 {
		apiResponse.Entries = &entriesResponse
	}

	// Update resource state with updated items
	configDataModel := operatorPolicyResourceModel{}
	err = stacuity.ConvertFromAPI(apiResponse, &configDataModel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting to TF",
			"Could not convert operator policy, unexpected error: "+err.Error(),
		)
		return
	}

	//reset
	if apiResponse.Entries != nil {
		configDataModel.Entries = &[]operatorPolicyEntryResource{}
		for _, entry := range *apiResponse.Entries {
			newEntry := operatorPolicyEntryResource{}
			newEntry.Iso3 = types.StringPointerValue(entry.Iso3)
			newEntry.OperatorId = types.Int32PointerValue(entry.OperatorId)
			newEntry.SteeringProfileEntryAction = types.StringValue(entry.SteeringProfileEntryAction.Moniker)
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
func (r *operatorPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state operatorPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing operator policy
	result, err := r.client.DeleteOperatorPolicy(state.Moniker.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting operator policy",
			"Could not delete operator policy, unexpected error: "+err.Error(),
		)
		return
	}

	if !result.Success {
		resp.Diagnostics.AddError(
			"Error Deleting operator policy",
			"Could not delete operator policy, unexpected errors: "+strings.Join(result.Messages, " "),
		)
		return
	}
}
