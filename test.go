package provider

import (
	"context"
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

func (r MinAPIHttpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	type State struct {
		ID           string
		ResponseBody string
	}

	state := State{
		ID:           url,
		ResponseBody: string(responseBody),
	}

	resp.State.Set(ctx, state)
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
	data := "g"
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
}

