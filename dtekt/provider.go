package dtekt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DTEKT_API_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dtekt_test":                resourceTest(),
			"dtekt_account":             resourceAccount(),
			"dtekt_uptime_monitor":      resourceUptimeMonitor(),
			"dtekt_alert_handler_slack": resourceAlertHandlerSlack(),
			"dtekt_performance_monitor": resourcePerformanceMonitor(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dtekt_tests": dataSourceTests(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	api_token := d.Get("api_token").(string)
	c, _ := NewClient(nil, &api_token)

	return c, diags
}
