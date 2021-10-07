package statuscake

import (
	"context"
	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceStatusCakeContactGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStatusCakeContactGroupCreate,
		ReadContext:   resourceStatusCakeContactGroupRead,
		UpdateContext: resourceStatusCakeContactGroupUpdate,
		DeleteContext: resourceStatusCakeContactGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the contact group",
			},
			"ping_url": {
				Type:        schema.TypeString, /* <uri> */
				Optional:    true,
				Description: "URL or IP address of an endpoint to push uptime events. Currently this only supports HTTP GET endpoints",
			},
			"email_addresses": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of email addresses",
			},
			"mobile_numbers": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of international format mobile phone numbers",
			},
			"integrations": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of integration IDs",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStatusCakeContactGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	req := client.CreateContactGroup(context.TODO()).
		Name(d.Get("name").(string))

	if v, ok := d.GetOk("ping_url"); ok {
		req = req.PingURL(v.(string))
	}
	if v, ok := d.GetOk("email_addresses"); ok {
		req = req.EmailAddresses(asListOfStrings(v))
	}
	if v, ok := d.GetOk("mobile_numbers"); ok {
		req = req.MobileNumbers(asListOfStrings(v))
	}
	if v, ok := d.GetOk("integrations"); ok {
		req = req.Integrations(asListOfStrings(v))
	}

	res, err := req.Execute()

	if err != nil {
		logStatusCakeAPIError(err)

		return asDiag(err.(statuscake.APIError))
	}

	logResponse(res)

	d.SetId(res.Data.NewID)

	return resourceStatusCakeContactGroupRead(ctx, d, meta)
}

func resourceStatusCakeContactGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	var diags diag.Diagnostics

	res, err := client.GetContactGroup(context.TODO(), d.Id()).Execute()

	if err != nil {
		logStatusCakeAPIError(err)

		if err.(statuscake.APIError).Status != 404 {
			return diag.FromErr(err)
		}
	}

	logResponse(res)

	if err := d.Set("name", res.Data.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ping_url", res.Data.PingURL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_addresses", res.Data.EmailAddresses); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mobile_numbers", res.Data.MobileNumbers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("integrations", res.Data.Integrations); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Data.ID)

	return diags
}

func resourceStatusCakeContactGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	if d.HasChangesExcept() {
		req := client.UpdateContactGroup(context.TODO(), d.Id())

		if d.HasChange("name") {
			req = req.Name(d.Get("name").(string))
		}
		if d.HasChange("ping_url") {
			req = req.PingURL(d.Get("ping_url").(string))
		}
		if d.HasChange("email_addresses") {
			req = req.EmailAddresses(asListOfStrings(d.Get("email_addresses")))
		}
		if d.HasChange("mobile_numbers") {
			req = req.MobileNumbers(asListOfStrings(d.Get("mobile_numbers")))
		}
		if d.HasChange("integrations") {
			req = req.Integrations(asListOfStrings(d.Get("integrations")))
		}

		if err := req.Execute(); err != nil {
			logStatusCakeAPIError(err)

			return asDiag(err.(statuscake.APIError))
		}
	}

	return resourceStatusCakeContactGroupRead(ctx, d, meta)
}

func resourceStatusCakeContactGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	var diags diag.Diagnostics

	err := client.DeleteContactGroup(context.TODO(), d.Id()).Execute()

	if err != nil {
		logStatusCakeAPIError(err)

		if err.(statuscake.APIError).Status != 404 {
			return diag.FromErr(err)
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Contact group has already been deleted",
		})
	}

	return diags
}
