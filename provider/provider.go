package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a Terraform ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"todo_task": resourceTask(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Retrieve the API URL from the provider configuration.
	url := d.Get("url").(string)

	// Normally, you would establish a connection to your backend service here.
	return &ProviderConfig{URL: url}, nil
}

// ProviderConfig holds the configuration details for the provider.
type ProviderConfig struct {
	URL string
}

// Resource definitions and CRUD functions go here...
