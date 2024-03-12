package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"reflect"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
    StatusCode    int
}

func (r *MinAPIHttpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "minapi_http"
}

func (r *MinAPIHttpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
The ` + "`minapi`" + ` resource makes an HTTP POST request to the given URL and exports
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
			"request_headers": schema.MapAttribute{
				Description: "A map of request header field names and values.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"status_code": schema.Int64Attribute{
                Description: "The HTTP status code of the response.",
                Computed:    true,
            },
		},
	}
}

func (r MinAPIHttpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// 1. Retrieve the configuration data from the request.
	var config MinAPIHttpResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. Extract relevant data from the configuration.
	url := config.Url.ValueString()
	method := config.Method.ValueString()
	payload := config.Payload.ValueString()
	requestHeaders := config.RequestHeaders

	// 3. Create an HTTP client and construct the request.
	client := &http.Client{}

	httpReq, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		resp.Diagnostics.AddError("Error creating HTTP request", err.Error())
		return
	}

	// 4. Add request headers.
	for key, value := range requestHeaders {
		httpReq.Header.Add(key, value)
	}

	// 5. Send the HTTP request and handle the response.
	response, err := client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Error making HTTP request", err.Error())
		return
	}

	// 6. Read the response body.
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		resp.Diagnostics.AddError("Error reading response body", err.Error())
		return
	}
	defer response.Body.Close()

    // Check the status code and add an error with response body if it's not 200
    if response.StatusCode != 200 {
        resp.Diagnostics.AddError(
            "Non-OK HTTP Status",
            fmt.Sprintf("HTTP request returned non-OK status code %d. Response body: %s", response.StatusCode, string(responseBody)),
        )
        return
    }

	// 8. Store the response data in the configuration.
	config.ID = types.StringValue(url)
	config.ResponseBody = types.StringValue(string(responseBody))
	config.StatusCode = types.Int64Value(int64(response.StatusCode))

	// 9. Set the updated configuration in the response state.
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MinAPIHttpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *MinAPIHttpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get the new configuration from the plan
    var plan MinAPIHttpResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Get the previous state
    var state MinAPIHttpResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Check if any input parameter has changed
    if plan.Url.ValueString() != state.Url.ValueString() ||
        plan.Method.ValueString() != state.Method.ValueString() ||
        plan.Payload.ValueString() != state.Payload.ValueString() ||
        !reflect.DeepEqual(plan.RequestHeaders, state.RequestHeaders) {

        // Remove the previous state
        resp.State.RemoveResource(ctx)

        // Perform the same logic as in the Create function
        url := plan.Url.ValueString()
        method := plan.Method.ValueString()
        payload := plan.Payload.ValueString()
        requestHeaders := plan.RequestHeaders

        client := &http.Client{}

        httpReq, err := http.NewRequest(method, url, strings.NewReader(payload))
        if err != nil {
            resp.Diagnostics.AddError("Error creating HTTP request", err.Error())
            return
        }

        for key, value := range requestHeaders {
            httpReq.Header.Add(key, value)
        }

        response, err := client.Do(httpReq)
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

        // Check the status code and add an error with response body if it's not 200
        if response.StatusCode != 200 {
            resp.Diagnostics.AddError(
                "Non-OK HTTP Status",
                fmt.Sprintf("HTTP request returned non-OK status code %d. Response body: %s", response.StatusCode, string(responseBody)),
            )
            return
        }

        plan.ID = types.StringValue(url)
        plan.ResponseBody = types.StringValue(string(responseBody))
        plan.StatusCode = types.Int64Value(int64(response.StatusCode))

        // Save the new state
        resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
    }
}

type MinAPIHttpResourceModel struct {
    Url            types.String `tfsdk:"url"`
    Payload        types.String `tfsdk:"payload"`
    Method         types.String `tfsdk:"method"`
    ID             types.String `tfsdk:"id"`
    ResponseBody   types.String `tfsdk:"response_body"`
    RequestHeaders map[string]string `tfsdk:"request_headers"`
    StatusCode     types.Int64 `tfsdk:"status_code"`
}