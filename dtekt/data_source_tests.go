package dtekt

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTests() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTestsRead,
		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"runs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"location": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTestsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	account_id := d.Get("account_id").(string)
	rc, err := c.GetTests(account_id)

	if err != nil {
		return diag.FromErr(err)
	}

	runs := normalizeTests(&rc.RunConfiguration.Runs)

	if err := d.Set("runs", runs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func normalizeTests(runs *[]Run) []interface{} {
	normalized := make([]interface{}, len(*runs), len(*runs))

	for i, run := range *runs {
		r := make(map[string]interface{})

		r["uuid"] = run.UUID
		r["url"] = run.Url
		r["schedule"] = run.Schedule
		r["location"] = run.Location

		normalized[i] = r
	}

	return normalized
}
