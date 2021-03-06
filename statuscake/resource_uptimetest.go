package statuscake

import (
	"context"
	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceStatusCakeUptimeTest() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a StatusCake Uptime Test",
		CreateContext: resourceStatusCakeUptimeTestCreate,
		ReadContext:   resourceStatusCakeUptimeTestRead,
		UpdateContext: resourceStatusCakeUptimeTestUpdate,
		DeleteContext: resourceStatusCakeUptimeTestDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the test",
			},
			// todo: ensure consistent casing?
			// todo: validate using UptimeTestType const
			// todo: include valid values in description
			"test_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Uptime test type",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"DNS",
						"HEAD",
						"HTTP",
						"PING",
						"SSH",
						"TCP",
					},
					false,
				),
			},
			"website_url": {
				Type:        schema.TypeString, /* <uri> */
				Required:    true,
				ForceNew:    true,
				Description: "URL or IP address of the website under test",
			},
			// todo: include valid values in description
			"check_rate": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of seconds between tests",
				ValidateFunc: validation.IntInSlice(
					[]int{
						0,
						30,
						60,
						300,
						900,
						1800,
						3600,
						86400,
					},
				),
			},
			"basic_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				ForceNew:    true,
				Description: "Basic authentication username",
			},
			"basic_pass": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Basic authentication password",
			},
			// todo: mention max value in description
			"confirmation": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "Number of confirmation servers to confirm downtime before an alert is triggered",
			},
			"contact_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of contact group IDs",
			},
			"custom_header": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON object. Represents headers to be sent when making requests",
			},
			"do_not_find": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to consider the test as down if the string in FindString is present within the response",
			},
			// todo: rename to "dns_ip"
			// todo: change to TypeList
			"dns_ip_csv": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comma separated list of IP addresses to compare against returned DNS records",
			},
			"dns_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Hostname or IP address of the nameserver to query",
			},
			"enable_ssl_alert": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Send an alert if the SSL certificate is soon to expire",
			},
			"final_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify where the redirect chain should end",
			},
			"find_string": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "String to look for within the response. Considered down if not found",
			},
			"follow_redirects": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow tests to follow redirects",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the hosting provider",
			},
			"include_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include header content in string match search",
			},
			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the test should be run",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Destination port for TCP tests",
			},
			"post_body": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "JSON object. This is converted to form data on request",
			},
			"post_raw": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Raw HTTP POST string to send to the server",
			},
			"regions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of regions on which to run tests. The values required for this parameter can be retrieved from the GET /v1/uptime-locations endpoint.",
			},
			"status_codes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of status codes that trigger an alert",
			},
			"tags": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of tags",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     40,
				Description: "How long to wait to receive the first byte",
			},
			"trigger_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4,
				Description: "The number of minutes to wait before sending an alert",
			},
			"cookie_storage": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable cookie storage",
			},
			"user_agent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User agent to be used when making requests",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceStatusCakeUptimeTestCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	req := client.CreateUptimeTest(context.TODO()).
		Name(d.Get("name").(string)).
		TestType(statuscake.UptimeTestType(d.Get("test_type").(string))).
		WebsiteURL(d.Get("website_url").(string)).
		CheckRate(statuscake.UptimeTestCheckRate(d.Get("check_rate").(int)))

	if v, ok := d.GetOk("basic_user"); ok {
		req = req.BasicUser(v.(string))
	}
	if v, ok := d.GetOk("basic_pass"); ok {
		req = req.BasicPass(v.(string))
	}
	if v, ok := d.GetOk("confirmation"); ok {
		req = req.Confirmation(int32(v.(int)))
	}
	if v, ok := d.GetOk("contact_groups"); ok {
		req = req.ContactGroups(asListOfStrings(v))
	}
	if v, ok := d.GetOk("custom_header"); ok {
		req = req.CustomHeader(v.(string))
	}
	if v, ok := d.GetOk("do_not_find"); ok {
		req = req.DoNotFind(v.(bool))
	}
	// todo: 'dns_ip_csv'
	if v, ok := d.GetOk("dns_server"); ok {
		req = req.DNSServer(v.(string))
	}
	if v, ok := d.GetOk("enable_ssl_alert"); ok {
		req = req.EnableSSLAlert(v.(bool))
	}
	if v, ok := d.GetOk("final_endpoint"); ok {
		req = req.FinalEndpoint(v.(string))
	}
	if v, ok := d.GetOk("find_string"); ok {
		req = req.FindString(v.(string))
	}
	if v, ok := d.GetOk("follow_redirects"); ok {
		req = req.FollowRedirects(v.(bool))
	}
	if v, ok := d.GetOk("host"); ok {
		req = req.Host(v.(string))
	}
	// todo: 'include_header'
	if v, ok := d.GetOk("paused"); ok {
		req = req.Paused(v.(bool))
	}
	if v, ok := d.GetOk("port"); ok {
		req = req.Port(int32(v.(int)))
	}
	if v, ok := d.GetOk("post_body"); ok {
		req = req.PostBody(v.(string))
	}
	if v, ok := d.GetOk("post_raw"); ok {
		req = req.PostRaw(v.(string))
	}
	// todo: 'regions'
	if v, ok := d.GetOk("status_codes"); ok {
		req = req.StatusCodes(asListOfStrings(v))
	}
	if v, ok := d.GetOk("tags"); ok {
		req = req.Tags(asListOfStrings(v))
	}
	if v, ok := d.GetOk("timeout"); ok {
		req = req.Timeout(int32(v.(int)))
	}
	if v, ok := d.GetOk("trigger_rate"); ok {
		req = req.TriggerRate(int32(v.(int)))
	}
	if v, ok := d.GetOk("cookie_storage"); ok {
		req = req.UseJAR(v.(bool))
	}
	if v, ok := d.GetOk("user_agent"); ok {
		req = req.UserAgent(v.(string))
	}

	res, err := req.Execute()

	if err != nil {
		logStatusCakeAPIError(err)

		return apiErrorDiag(err)
	}

	logResponse(res)

	d.SetId(res.Data.NewID)

	// currently, there isn't any way to set status_codes to "nothing" when creating
	// a new uptime test (meaning it always has a default value of all the codes)
	// and terraform has no way of providing a default value for a list.
	//
	// so we have to do an update straight after creating the uptime test to ensure
	// that the state matches what terraform expects
	// todo: discuss with StatusCake if this could be supported somehow?
	err = client.UpdateUptimeTest(context.TODO(), d.Id()).
		StatusCodes(asListOfStrings(d.Get("status_codes"))).
		Execute()

	if err != nil {
		logStatusCakeAPIError(err)

		return apiErrorDiag(err)
	}

	logResponse(res)

	return resourceStatusCakeUptimeTestRead(ctx, d, meta)
}

func resourceStatusCakeUptimeTestRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	var diags diag.Diagnostics

	res, err := client.GetUptimeTest(context.TODO(), d.Id()).Execute()

	if err != nil {
		logStatusCakeAPIError(err)

		if !isNotFoundAPIError(err) {
			return diag.FromErr(err)
		}
	}

	logResponse(res)

	if err := d.Set("name", res.Data.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("test_type", res.Data.TestType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("website_url", res.Data.WebsiteURL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("check_rate", res.Data.CheckRate); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("confirmation", res.Data.Confirmation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("contact_groups", res.Data.ContactGroups); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("custom_header", res.Data.CustomHeader); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("do_not_find", res.Data.DoNotFind); err != nil {
		return diag.FromErr(err)
	}
	// todo: 'dns_ip_csv'
	if err := d.Set("dns_server", res.Data.DNSServer); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enable_ssl_alert", res.Data.EnableSSLAlert); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("final_endpoint", res.Data.FinalEndpoint); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("find_string", res.Data.FindString); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("follow_redirects", res.Data.FollowRedirects); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("host", res.Data.Host); err != nil {
		return diag.FromErr(err)
	}
	// todo: 'include_header'
	if err := d.Set("paused", res.Data.Paused); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", res.Data.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("post_body", res.Data.PostBody); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("post_raw", res.Data.PostRaw); err != nil {
		return diag.FromErr(err)
	}
	// todo: 'regions'
	if err := d.Set("status_codes", res.Data.StatusCodes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", res.Data.Tags); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("timeout", res.Data.Timeout); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("trigger_rate", res.Data.TriggerRate); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cookie_storage", res.Data.UseJAR); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_agent", res.Data.UserAgent); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(res.Data.ID)

	return diags
}

func resourceStatusCakeUptimeTestUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	if d.HasChangesExcept() {
		req := client.UpdateUptimeTest(context.TODO(), d.Id())

		if d.HasChange("name") {
			req = req.Name(d.Get("name").(string))
		}
		if d.HasChange("check_rate") {
			req = req.CheckRate(statuscake.UptimeTestCheckRate(d.Get("check_rate").(int)))
		}
		if d.HasChange("basic_user") {
			req = req.BasicUser(d.Get("basic_user").(string))
		}
		if d.HasChange("basic_pass") {
			req = req.BasicPass(d.Get("basic_pass").(string))
		}
		if d.HasChange("confirmation") {
			req = req.Confirmation(int32(d.Get("confirmation").(int)))
		}
		if d.HasChange("contact_groups") {
			req = req.ContactGroups(asListOfStrings(d.Get("contact_groups")))
		}
		if d.HasChange("custom_header") {
			req = req.CustomHeader(d.Get("custom_header").(string))
		}
		if d.HasChange("do_not_find") {
			req = req.DoNotFind(d.Get("do_not_find").(bool))
		}
		// todo: 'dns_ip_csv'
		if d.HasChange("dns_server") {
			req = req.DNSServer(d.Get("dns_server").(string))
		}
		if d.HasChange("enable_ssl_alert") {
			req = req.EnableSSLAlert(d.Get("enable_ssl_alert").(bool))
		}
		if d.HasChange("final_endpoint") {
			req = req.FinalEndpoint(d.Get("final_endpoint").(string))
		}
		if d.HasChange("find_string") {
			req = req.FindString(d.Get("find_string").(string))
		}
		if d.HasChange("follow_redirects") {
			req = req.FollowRedirects(d.Get("follow_redirects").(bool))
		}
		if d.HasChange("host") {
			req = req.Host(d.Get("host").(string))
		}
		// todo: 'include_header'
		if d.HasChange("paused") {
			req = req.Paused(d.Get("paused").(bool))
		}
		if d.HasChange("port") {
			req = req.Port(int32(d.Get("port").(int)))
		}
		if d.HasChange("post_body") {
			req = req.PostBody(d.Get("post_body").(string))
		}
		if d.HasChange("post_raw") {
			req = req.PostRaw(d.Get("post_raw").(string))
		}
		// todo: 'regions'
		if d.HasChange("status_codes") {
			req = req.StatusCodes(asListOfStrings(d.Get("status_codes")))
		}
		if d.HasChange("tags") {
			req = req.Tags(asListOfStrings(d.Get("tags")))
		}
		if d.HasChange("timeout") {
			req = req.Timeout(int32(d.Get("timeout").(int)))
		}
		if d.HasChange("trigger_rate") {
			req = req.TriggerRate(int32(d.Get("trigger_rate").(int)))
		}
		if d.HasChange("cookie_storage") {
			req = req.UseJAR(d.Get("cookie_storage").(bool))
		}
		if d.HasChange("user_agent") {
			req = req.UserAgent(d.Get("user_agent").(string))
		}

		if err := req.Execute(); err != nil {
			logStatusCakeAPIError(err)

			return apiErrorDiag(err)
		}
	}

	return resourceStatusCakeUptimeTestRead(ctx, d, meta)
}

func resourceStatusCakeUptimeTestDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*statuscake.APIClient)

	var diags diag.Diagnostics

	err := client.DeleteUptimeTest(context.TODO(), d.Id()).Execute()

	if err != nil {
		logStatusCakeAPIError(err)

		if !isNotFoundAPIError(err) {
			return diag.FromErr(err)
		}

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Uptime test has already been deleted",
		})
	}

	return diags
}
