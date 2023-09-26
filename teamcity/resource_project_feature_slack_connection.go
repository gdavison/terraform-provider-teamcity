package teamcity

import (
	"context"
	"fmt"
	"log"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectFeatureSlackConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectFeatureSlackConnectionCreate,
		Read:   resourceProjectFeatureSlackConnectionRead,
		Update: resourceProjectFeatureSlackConnectionUpdate,
		Delete: resourceProjectFeatureSlackConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceProjectFeatureSlackConnectionImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceProjectFeatureSlackConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Get("project_id").(string)
	service := client.ProjectFeatureService(projectId)

	feature := api.NewProjectFeatureSlackConnection(projectId, api.ProjectFeatureSlackConnectionOptions{
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		DisplayName:  d.Get("display_name").(string),
		Token:        d.Get("token").(string),
	})

	created, err := service.Create(feature)
	if err != nil {
		return err
	}

	d.SetId(created.ID())

	return resourceProjectFeatureSlackConnectionRead(d, meta)
}

func resourceProjectFeatureSlackConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Get("project_id").(string)
	featureId := d.Id()

	service := client.ProjectFeatureService(projectId)

	slackConnection := api.NewProjectFeatureSlackConnection(projectId, api.ProjectFeatureSlackConnectionOptions{
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		DisplayName:  d.Get("display_name").(string),
		Token:        d.Get("token").(string),
	})
	slackConnection.SetID(featureId)

	if _, err := service.Update(slackConnection); err != nil {
		return err
	}

	return resourceProjectFeatureSlackConnectionRead(d, meta)
}

func resourceProjectFeatureSlackConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Get("project_id").(string)
	featureId := d.Id()

	service := client.ProjectFeatureService(projectId)
	feature, err := service.GetByID(featureId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[DEBUG] Project Feature Slack Notifier (%s) not found, removing from state", featureId)
			d.SetId("")
			return nil
		}

		return err
	}

	slackConnection, ok := feature.(*api.ProjectFeatureSlackConnection)
	if !ok {
		return fmt.Errorf("Expected a VersionedSettings Feature but wasn't!")
	}

	d.Set("client_id", string(slackConnection.Options.ClientId))
	d.Set("client_secret", d.Get("client_secret"))
	d.Set("display_name", string(slackConnection.Options.DisplayName))
	d.Set("token", d.Get("token"))

	return nil
}

func resourceProjectFeatureSlackConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Get("project_id").(string)
	featureId := d.Id()

	service := client.ProjectFeatureService(projectId)
	feature, err := service.GetByID(featureId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}

		return err
	}

	return service.Delete(feature.ID())
}

func resourceProjectFeatureSlackConnectionImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return nil, fmt.Errorf("unexpected format for ID (%s), use: '<project-id>/<feature-id>'", d.Id())
	}

	projectId := idParts[0]
	featureId := idParts[1]

	d.Set("project_id", projectId)
	d.SetId(featureId)

	return []*schema.ResourceData{d}, nil

}
