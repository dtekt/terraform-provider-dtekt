package dtekt

import (
	"context"
	// "strconv"

	// "fmt"
	// "errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountCreate,
		UpdateContext: resourceAccountUpdate,
		ReadContext:   resourceAccountRead,
		DeleteContext: resourceAccountDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	newAccount := Account{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	resp, err := c.CreateAccount(newAccount)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.Id)

	resourceAccountRead(ctx, d, m)

	return diags
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	accountId := d.Id()

	_, err := c.GetAccount(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceAccountRead(ctx, d, m)
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)

	var diags diag.Diagnostics

	accountId := d.Id()

	err := c.DeleteAccount(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
