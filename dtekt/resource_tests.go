package dtekt

import (
	"context"

	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTestCreate,
		UpdateContext: resourceTestUpdate,
		ReadContext:   resourceTestRead,
		DeleteContext: resourceTestDelete,
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
			"schedule": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTestCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	newTest := Run{
		Url:      d.Get("url").(string),
		Schedule: d.Get("schedule").(int),
		Location: d.Get("location").(string),
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Warning Message Summary",
		Detail:   fmt.Sprintf("Received: URL=%s Schedule=%d Location=%s", newTest.Url, newTest.Schedule, newTest.Location),
	})

	account_id := d.Get("account_id").(string)

	resp, err := c.CreateTest(newTest, account_id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.UUID)

	resourceTestRead(ctx, d, m)

	return diags
}

func resourceTestRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	testUuid := d.Id()
	account_id := d.Get("account_id").(string)

	_, err := c.GetTest(testUuid, account_id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceTestUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTestRead(ctx, d, m)
}

func resourceTestDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	testUuid := d.Id()
	account_id := d.Get("account_id").(string)

	err := c.DeleteTest(testUuid, account_id)

	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
