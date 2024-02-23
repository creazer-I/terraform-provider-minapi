package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateVerb(val interface{}, key string) (warns []string, errs []error) {
	// Validation logic here
	return
}

func resourceMinApi() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMinApiCreate,
		ReadContext:   nil, // Read operation integrated into Create
		UpdateContext: resourceMinApiUpdate,
		DeleteContext: resourceMinApiDelete,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "GET",
			},
			"request_headers": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			// Define other attributes as needed
		},
	}
}

func resourceMinApiCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	url := d.Get("url").(string)
	method := d.Get("method").(string)
	headers := d.Get("request_headers").(map[string]interface{})

	// Perform the HTTP request
	client := &http.Client{
		Timeout: time.Second * 10, // Example timeout, adjust as needed
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Read response body
	// Update resource data with response attributes
	// SetId as needed
	d.SetId(url)

	return nil
}

func resourceMinApiUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Update logic if necessary
	return nil
}

func resourceMinApiDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Delete logic if necessary
	return nil
}
