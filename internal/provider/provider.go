// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure FunctionsProvider satisfies various provider interfaces.
var _ provider.ProviderWithFunctions = &FunctionsProvider{}

// FunctionsProvider defines the provider implementation.
type FunctionsProvider struct {
	version string
}

// FunctionsProviderModel describes the provider data model.
type FunctionsProviderModel struct{}

func (p *FunctionsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "functions"
	resp.Version = p.version
}

func (p *FunctionsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *FunctionsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *FunctionsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *FunctionsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *FunctionsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewBase64TarGzFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FunctionsProvider{
			version: version,
		}
	}
}
