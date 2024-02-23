package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*MinAPIHttpResource)(nil)

func NewMinAPIHttpResource() resource.Resource {
	return &MinAPIHttpResource{}
}

type MinAPIHttpResource struct{}

func (r *MinAPIHttpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "http"
}

func (r *MinAPIHttpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
The ` + "`http`" + ` resource makes an HTTP POST request to the given URL and exports
information about the response.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The URL used for the request.",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The URL for the request.",
				Required:    true,
			},
			"payload": schema.StringAttribute{
				Description: "The payload to be sent in the POST request.",
				Required:    true,
			},
			"response_body": schema.StringAttribute{
				Description: "The response body returned as a string.",
				Computed:    true,
			},
		},
	}
}

func (r *MinAPIHttpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config MinAPIHttpResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := config.Url.ValueString()
	payload := config.Payload.ValueString()

	response, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		resp.Diagnostics.AddError("Error making HTTP request", err.Error())
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		resp.Diagnostics.AddError("Error reading response body", err.Error())
		return
	}

	resp.State.Set("id", url)
	resp.State.Set("response_body", types.StringValue(string(responseBody)))
}

type MinAPIHttpResourceModel struct {
	Url     types.String `tfsdk:"url"`
	Payload types.String `tfsdk:"payload"`
}
