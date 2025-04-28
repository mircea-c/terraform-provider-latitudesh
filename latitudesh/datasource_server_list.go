package latitudesh

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/latitudesh/latitudesh-go"
)

func dataSourceServerList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServerListRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Description: "The project identifier.",
				Required:    true,
			},
			"plan_name": {
				Type:        schema.TypeString,
				Description: "The plan name.",
				Optional:    true,
			},
			"locations": {
				Type:        schema.TypeMap,
				Description: "The list of servers belonging to the project.",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeMap,
				Description: "The list of servers belonging to the project.",
				Computed:    true,
			},
		},
	}
}

func dataSourceServerListRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*api.Client)

	options := new(api.ListOptions)
	options = options.Filter("plan", d.Get("plan_name").(string))
	options = options.AddParam("page[size]", "1000")

	project := d.Get("project").(string)
	servers, _, err := c.Servers.List(project, options)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(servers) == 0 {
		return diag.Errorf("No servers found for project ID %s", project)
	}

	locations := map[string]string{}
	for _, server := range servers {
		locations[server.ID] = server.Region.Site.Slug
	}

	if len(locations) == 0 {
		return diag.Errorf("No servers found for project ID %s", project)
	}

	labels := map[string]string{}
	for _, server := range servers {
		labels[server.ID] = server.Label
	}

	if len(labels) == 0 {
		return diag.Errorf("No servers found for project ID %s", project)
	}

	d.SetId(d.Get("plan_name").(string))

	if err = d.Set("locations", locations); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("labels", labels); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
