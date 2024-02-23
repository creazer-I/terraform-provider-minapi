package provider

import (
	"context"
	//"encoding/json"
	"io"
	"net/http"

	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	//"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ resource.Resource = (*MinAPIHttpResource)(nil)

func NewMinAPIHttpResource() resource.Resource {
	return &MinAPIHttpResource{}
}

type MinAPIHttpResource struct{}

type State struct {
	Url           string
	Payload       string
	Method        string
	ID            string
	ResponseBody  string
}

func (r *MinAPIHttpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "minapi_http"
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
				Optional: true,

			},
			"method": schema.StringAttribute{
				Description: "The HTTP method for the request.",
				Required:    true,
			},
			"response_body": schema.StringAttribute{
				Description: "The response body returned as a string.",
				Computed:    true,
			},
		},
	}
}

func (r MinAPIHttpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var config MinAPIHttpResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := config.Url.ValueString()
	method := config.Method.ValueString()
	payload := config.Payload.ValueString()

	client := &http.Client{}

	// if method == "GET" {
	// }

	httpReq, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		resp.Diagnostics.AddError("Error creating HTTP request", err.Error())
		return
	}

	response, err := client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error making HTTP request", err.Error())
		return
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		resp.Diagnostics.AddError("Error reading response body", err.Error())
		return
	}
	defer response.Body.Close()


	config.ID = types.StringValue(url)
	config.ResponseBody = types.StringValue(string(responseBody))
	diags = resp.State.Set(ctx, config)
	if diags.HasError() {
		for _, diag := range diags {
			resp.Diagnostics.AddError("Error setting state", diag.Summary())
		}
		return
	}
}

func (r *MinAPIHttpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *MinAPIHttpResourceModel

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

func (r *MinAPIHttpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
}

func (r *MinAPIHttpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *MinAPIHttpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update MinApi, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}


type MinAPIHttpResourceModel struct {
	Url     types.String `tfsdk:"url"`
	Payload types.String `tfsdk:"payload"`
	Method  types.String `tfsdk:"method"`
	ID          types.String `tfsdk:"id"`
	ResponseBody types.String `tfsdk:"response_body"`
}