package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heroku/docker-registry-client/registry"
)

const defaultRetryAttempts = 3
const defaultRetryDelay = 10

func dataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,

		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_EMERGE_REGISTRY", ""),
			},

			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_EMERGE_REPOSITORY", ""),
			},

			"tag": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_EMERGE_TAG", ""),
			},

			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_EMERGE_USERNAME", ""),
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DOCKER_EMERGE_PASSWORD", ""),
				Sensitive:   true,
			},

			"retry": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attempts": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  defaultRetryAttempts,
						},
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  defaultRetryDelay,
						},
					},
				},
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	url := d.Get("registry").(string)
	repository := d.Get("repository").(string)
	tag := d.Get("tag").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	hub, err := registry.New(url, username, password)
	if err != nil {
		return append(diags, diag.Errorf("Error creating hub: %s", err)...)
	}

	if v, ok := d.GetOk("retry"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		retry := v.([]interface{})[0].(map[string]interface{})
		attempts := retry["attempts"].(int)
		wait := time.Duration(retry["delay"].(int)) * time.Second
		count := 0

		for count < attempts {
			tags, err := hub.Tags(repository)
			if err != nil {
				return append(diags, diag.Errorf("Error getting tags for repository: %s", err)...)
			}

			for _, v := range tags {
				if v == tag {
					// Tag exists
					d.SetId(fmt.Sprintf("%s:%s", repository, tag))
					return diags
				}
			}

			time.Sleep(wait)
			count++
		}
	}

	return append(diags, diag.Errorf("Docker image does not exist")...)
}
