// SPDX-License-Identifier: MPL-2.0
package provider

import (
	"context"
	"os"

	stacuity "stacuity.com/go_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure StacuityProvider satisfies various provider interfaces.
var _ provider.Provider = &StacuityProvider{}

// StacuityProvider defines the provider implementation.
type StacuityProvider struct {
	version string
}

// StacuityProviderModel describes the provider data model.
type StacuityProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *StacuityProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "stacuity"
	resp.Version = p.version
}

func (p *StacuityProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Stacuity.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URL for Stacuity API. May also be provided via STACUITY_HOST environment variable. Optional",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "Token for Stacuity API. May also be provided via STACUITY_TOKEN environment variable.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

type stacuityProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func (p *StacuityProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring Stacuity client")
	// Retrieve provider data from configuration
	var config stacuityProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("STACUITY_HOST")
	token := os.Getenv("STACUITY_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if host == "" {
		host = "https://api.stacuity.com/api/v1"
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Stacuity API Host",
			"The provider cannot create the Stacuity API client as there is a missing or empty value for the Stacuity API host. "+
				"Set the host value in the configuration or use the STACUITY_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Stacuity API Token",
			"The provider cannot create the Stacuity API client as there is a missing or empty value for the Stacuity API token. "+
				"Set the token value in the configuration or use the STACUITY_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Stacuity client")

	// Create a new Stacuity client using the configuration values
	client, err := stacuity.NewClient(&host, &token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Stacuity API Client",
			"An unexpected error occurred when creating the Stacuity API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Stacuity Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Stacuity client")
}

func (p *StacuityProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVSliceResource, NewRoutingTargetResource, NewRoutingPolicyResource, NewEndpointGroupResource,
		NewEventMapResource, NewEventHandlerResource, NewOperatorPolicyResource, NewRegionalPolicyResource,
	}
}

func VSliceDataSource() datasource.DataSource {
	return &vSlicesDataSource{}
}

func RoutingTargetDataSource() datasource.DataSource {
	return &routingTargetDataSource{}
}

func RoutingPolicyDataSource() datasource.DataSource {
	return &routingPolicyDataSource{}
}

func EndpointGroupDataSource() datasource.DataSource {
	return &endpointGroupDataSource{}
}

func EventMapDataSource() datasource.DataSource {
	return &eventMapDataSource{}
}

func EventHandlerDataSource() datasource.DataSource {
	return &eventHandlerDataSource{}
}

func OperatorPolicyDataSource() datasource.DataSource {
	return &operatorPolicyDataSource{}
}

func RegionalPolicyDataSource() datasource.DataSource {
	return &regionalPolicyDataSource{}
}

// DataSources defines the data sources implemented in the provider.
func (p *StacuityProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		VSliceDataSource, RoutingTargetDataSource, RoutingPolicyDataSource, EndpointGroupDataSource, EventMapDataSource,
		EventHandlerDataSource, OperatorPolicyDataSource, RegionalPolicyDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &StacuityProvider{
			version: version,
		}
	}
}
