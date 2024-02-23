// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"


	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/net/http/httpproxy"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

)

func validateVerb(val interface{}, key string) (warns []string, errs []error) {
	if v, ok := val.(string); ok {
		switch v {
		case
			http.MethodGet,
			http.MethodPost,
			http.MethodHead,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodPut:
			break
		default:
			errs = append(errs, fmt.Errorf("%s must be GET|POST|HEAD|DELETE|PATCH|PUT, got: %s", key, v))
		}
	} else {
		errs = append(errs, fmt.Errorf("error parsing method"))
	}
	return
}

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MinApiResource{}
var _ resource.ResourceWithImportState = &MinApiResource{}

func NewMinApiResource() resource.Resource {
	return &NewMinApiResource{}
}

// MinApiResource defines the resource implementation.
type MinApiResource struct {
	client *minapi.Client
}

// MinApiResourceModel describes the resource data model.
// type MinApiResourceModel struct {
// 	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
// 	Defaulted             types.String `tfsdk:"defaulted"`
// 	Id                    types.String `tfsdk:"id"`
// }

func (r *MinApiResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_MinApi"
}

func (r *MinApiResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "MinApi resource",

		Attributes: map[string]schema.Attribute{
			"configurable_attribute": schema.StringAttribute{
				MarkdownDescription: "MinApi configurable attribute",
				Optional:            true,
			},
			"defaulted": schema.StringAttribute{
				MarkdownDescription: "MinApi configurable attribute with default value",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("MinApi value when not configured"),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "MinApi identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *MinApiResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *MinApiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MinApiResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create MinApi, got error: %s", err))
	//     return
	// }

	// For the purposes of this MinApi code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue("MinApi-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MinApiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MinApiResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read MinApi, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MinApiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *MinApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MinApiResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete MinApi, got error: %s", err))
	//     return
	// }
}

func (r *MinApiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
