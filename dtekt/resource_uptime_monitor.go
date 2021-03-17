package dtekt

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUptimeMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUptimeMonitorCreate,
		UpdateContext: resourceUptimeMonitorUpdate,
		ReadContext:   resourceUptimeMonitorRead,
		DeleteContext: resourceUptimeMonitorDelete,
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"match_string": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"success_code": &schema.Schema{
				Type:     schema.TypeInt,
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
				Default:      "1m",
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
			"header": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"alert": AlertSchema,
		},
	}
}

func buildUptimeMonitor(d *schema.ResourceData) UptimeMonitor {
	// Convert headers scheme.Set to map which will be
	// JSON-marshalled
	headers := d.Get("header").(*schema.Set)
	headersMap := make(map[string]string)
	for _, h := range headers.List() {
		header := h.(map[string]interface{})
		headersMap[header["name"].(string)] = header["value"].(string)
	}

	// Build locationConfig
	locations := d.Get("locations").(*schema.Set).List()
	locationConfig := BuildLocationConfig(locations)

	// if DevEnv() {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Summary:  "Built locationConfig:",
	// 		Detail:   fmt.Sprintf("%+v", locationConfig),
	// 	})
	// }

	return UptimeMonitor{
		Url:            d.Get("url").(string),
		MatchString:    d.Get("match_string").(string),
		SuccessCode:    d.Get("success_code").(int),
		Every:          CovertEvery(d.Get("every")).(string),
		LocationConfig: locationConfig,
		Headers:        headersMap,
	}
}

func resourceUptimeMonitorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	accountId := d.Get("account_id").(string)

	newUptimeMonitor := buildUptimeMonitor(d)

	resp, err := c.CreateUptimeMonitor(newUptimeMonitor, accountId)
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
			Handlers:  alert["handlers"].(*schema.Set).List(),
			MonitorId: resp.UUID,
			Window:    alert["window"].(string),
		}, accountId)
	}

	d.SetId(resp.UUID)

	resourceUptimeMonitorRead(ctx, d, m)
	return diags
}

func resourceUptimeMonitorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	uptimeMonitorUuid := d.Id()
	accountId := d.Get("account_id").(string)

	resp, err := c.GetUptimeMonitor(uptimeMonitorUuid, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("sucess_code", resp.SuccessCode)
	d.Set("url", resp.Url)
	d.Set("every", resp.Every)

	if err := d.Set("alert", flattenAlertDefinitions(&resp.AlertDefinition)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("locations", FlattenLocationConfig(&resp.LocationConfig)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func flattenAlertDefinitions(alertDefinitions *[]AlertDefinition) []interface{} {
	if alertDefinitions != nil {
		alerts := make([]interface{}, len(*alertDefinitions), len(*alertDefinitions))
		for i, alertDefinition := range *alertDefinitions {
			ad := make(map[string]interface{})

			var handlers []string
			for _, handler := range alertDefinition.Handlers {
				h := handler.(map[string]interface{})
				handlers = append(handlers, h["uuid"].(string))
			}

			ad["warn"] = alertDefinition.Warn
			ad["uuid"] = alertDefinition.UUID
			ad["crit"] = alertDefinition.Crit
			ad["metric"] = alertDefinition.Metric
			ad["handlers"] = handlers
			ad["window"] = alertDefinition.Window

			alerts[i] = ad
		}

		return alerts
	}

	return make([]interface{}, 0)
}

func resourceUptimeMonitorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	monitorId := d.Id()
	accountId := d.Get("account_id").(string)

	// Update entire object when one of the simple attributes changes
	if d.HasChange("locations") || d.HasChange("every") || d.HasChange("url") || d.HasChange("success_code") {
		updatedUptimeMonitor := buildUptimeMonitor(d)
		_, err := c.UpdateUptimeMonitor(updatedUptimeMonitor, accountId, monitorId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Updating alerts requires totally separate API calls that's why it is
	// performed in a seperate block
	if d.HasChange("alert") {
		old, new := d.GetChange("alert")

		o := old.(*schema.Set)
		// panic(fmt.Sprintf("%+v", o))

		n := new.(*schema.Set)
		// panic(fmt.Sprintf("%+v", n))

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

				_, err := c.CreateAlertDefinition(AlertDefinition{
					Warn:      alert["warn"].(float64),
					Crit:      alert["crit"].(float64),
					Metric:    alert["metric"].(string),
					Handlers:  alert["handlers"].(*schema.Set).List(),
					MonitorId: monitorId,
					Window:    alert["window"].(string),
				}, accountId)

				if err != nil {
					diags = append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Unable to update Alert Definition. Error:",
						Detail:   fmt.Sprintf("%+v", err),
					})
				}
			}
		}
	}

	resourceUptimeMonitorRead(ctx, d, m)
	return diags
}

func resourceUptimeMonitorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	uptimeMonitorUuid := d.Id()
	accountId := d.Get("account_id").(string)

	err := c.DeleteUptimeMonitor(uptimeMonitorUuid, accountId)

	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
