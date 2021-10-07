package statuscake

import (
	"context"
	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_token": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("STATUSCAKE_API_KEY", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"statuscake_contact_group": ResourceStatusCakeContactGroup(),
				"statuscake_uptime_test":   ResourceStatusCakeUptimeTest(),
			},
			DataSourcesMap: map[string]*schema.Resource{},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiToken := d.Get("api_token").(string)

		var diags diag.Diagnostics

		if apiToken == "" {
			return nil, diag.Errorf("Missing api_token")
		}

		client := statuscake.NewAPIClient(apiToken)

		return client, diags
	}
}
