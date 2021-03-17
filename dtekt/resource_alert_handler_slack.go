package dtekt

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlertHandlerSlack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertHandlerSlackCreate,
		UpdateContext: resourceAlertHandlerSlackUpdate,
		ReadContext:   resourceAlertHandlerSlackRead,
		DeleteContext: resourceAlertHandlerSlackDelete,
		Schema: map[string]*schema.Schema{
			"webhook_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlertHandlerSlackCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// diags = append(diags, diag.Diagnostic{
	// 	Severity: diag.Warning,
	// 	Summary:  "Warning Message Summary",
	// 	Detail:   fmt.Sprintf("%v \n %v", headersMap, headersMap2),
	// })

	opts := make(map[string]interface{})
	opts["webhook_url"] = d.Get("webhook_url").(string)
	newAlertHandlerSlack := AlertHandler{
		Kind:    "slack",
		Options: opts,
		Name:    d.Get("name").(string),
	}

	// newAlertHandlerSlack.Options =

	// diags = append(diags, diag.Diagnostic{
	// 	Severity: diag.Warning,
	// 	Summary:  "Warning Message Summary",
	// 	Detail:   fmt.Sprintf("Received: URL=%s Schedule=%d Location=%s", newTest.Url, newTest.Schedule, newTest.Location),
	// })

	account_id := d.Get("account_id").(string)

	resp, err := c.CreateAlertHandlerSlack(newAlertHandlerSlack, account_id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.UUID)

	resourceAlertHandlerSlackRead(ctx, d, m)

	return diags
}

func resourceAlertHandlerSlackRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	AlertHandlerSlackUuid := d.Id()
	accountId := d.Get("account_id").(string)

	resp, err := c.GetAlertHandlerSlack(AlertHandlerSlackUuid, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Warning Message Summary",
		Detail:   fmt.Sprintf("Resp %+v", resp),
	})

	d.Set("webhook_url", resp.Options["webhook_url"])

	return diags
}

func resourceAlertHandlerSlackUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	accountId := d.Get("account_id").(string)
	handlerId := d.Id()

	opts := make(map[string]interface{})
	opts["webhook_url"] = d.Get("webhook_url").(string)
	newAlertHandlerSlack := AlertHandler{
		Kind:    "slack",
		Options: opts,
		Name:    d.Get("name").(string),
	}

	_, err := c.UpdateAlertHandlerSlack(newAlertHandlerSlack, accountId, handlerId)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAlertHandlerSlackRead(ctx, d, m)
}

func resourceAlertHandlerSlackDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	AlertHandlerSlackUuid := d.Id()
	account_id := d.Get("account_id").(string)

	err := c.DeleteAlertHandlerSlack(AlertHandlerSlackUuid, account_id)

	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
