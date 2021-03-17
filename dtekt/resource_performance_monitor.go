package dtekt

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePerformanceMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePerformanceMonitorCreate,
		UpdateContext: resourcePerformanceMonitorUpdate,
		ReadContext:   resourcePerformanceMonitorRead,
		DeleteContext: resourcePerformanceMonitorDelete,
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"every": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "15m",
				StateFunc:    everyStateFunc,
				ValidateFunc: everyValidatorFunc,
			},
			"locations": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"alert": AlertSchema,
		},
	}
}

func resourcePerformanceMonitorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	accountId := d.Get("account_id").(string)

	// Convert headers scheme.Set to map which will be
	// JSON-marshalled

	locations := d.Get("locations").(*schema.Set).List()
	locationConfig := BuildLocationConfig(locations)

	newPerformanceMonitor := PerformanceMonitor{
		Url:            d.Get("url").(string),
		Every:          d.Get("every").(string),
		LocationConfig: locationConfig,
	}

	resp, err := c.CreatePerformanceMonitor(newPerformanceMonitor, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	alerts := d.Get("alert").(*schema.Set)
	for _, a := range alerts.List() {
		alert := a.(map[string]interface{})
		c.CreateAlertDefinition(AlertDefinition{
			Warn:      alert["warn"].(float64),
			Crit:      alert["crit"].(float64),
			Metric:    alert["metric"].(string),
			Handlers:  alert["handlers"].([]interface{}),
			MonitorId: resp.UUID,
			Window:    alert["window"].(string),
		}, accountId)
	}

	d.SetId(resp.UUID)

	resourcePerformanceMonitorRead(ctx, d, m)
	return diags
}

func resourcePerformanceMonitorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	PerformanceMonitorUuid := d.Id()
	accountId := d.Get("account_id").(string)

	resp, err := c.GetPerformanceMonitor(PerformanceMonitorUuid, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("alert", flattenAlertDefinitions(&resp.AlertDefinition)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", FlattenLocationConfig(&resp.LocationConfig)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("url", resp.Url)
	d.Set("every", resp.Every)

	return diags
}

func resourcePerformanceMonitorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	monitorId := d.Id()
	accountId := d.Get("account_id").(string)

	if d.HasChange("alert") {
		old, new := d.GetChange("alert")

		o := old.(*schema.Set)
		n := new.(*schema.Set)

		toDelete := o.Difference(n)
		toCreate := n.Difference(o)

		if toDelete.Len() > 0 {
			for _, a := range toDelete.List() {
				alert := a.(map[string]interface{})
				c.DeleteAlertDefinition(alert["uuid"].(string), accountId)
			}
		}

		if toCreate.Len() > 0 {
			for _, a := range toCreate.List() {
				alert := a.(map[string]interface{})

				c.CreateAlertDefinition(AlertDefinition{
					Warn:      alert["warn"].(float64),
					Crit:      alert["crit"].(float64),
					Metric:    alert["metric"].(string),
					Handlers:  alert["handlers"].([]interface{}),
					MonitorId: monitorId,
					Window:    alert["window"].(string),
				}, accountId)
			}
		}
	}

	return diags
}

func resourcePerformanceMonitorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	PerformanceMonitorUuid := d.Id()
	accountId := d.Get("account_id").(string)

	err := c.DeletePerformanceMonitor(PerformanceMonitorUuid, accountId)

	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
