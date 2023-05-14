package teamcity

import (
	"context"
	"fmt"
	"log"
	"strings"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectFeatureSlackNotifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectFeatureSlackNotifierCreate,
		Read:   resourceProjectFeatureSlackNotifierRead,
		Update: resourceProjectFeatureSlackNotifierUpdate,
		Delete: resourceProjectFeatureSlackNotifierDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceProjectFeatureSlackNotifierImport,
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

func resourceProjectFeatureSlackNotifierCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Get("project_id").(string)
	service := client.ProjectFeatureService(projectId)

	feature := api.NewProjectFeatureSlackNotifier(projectId, api.ProjectFeatureSlackNotifierOptions{
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

	return resourceProjectFeatureSlackNotifierRead(d, meta)
}

func resourceProjectFeatureSlackNotifierUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Id()

	service := client.ProjectFeatureService(projectId)
	feature, err := service.GetByType("versionedSettings")
	if err != nil {
		return err
	}

	vcsFeature, ok := feature.(*api.ProjectFeatureVersionedSettings)
	if !ok {
		return fmt.Errorf("Expected a VersionedSettings Feature but wasn't!")
	}

	if d.HasChange("build_settings") {
		vcsFeature.Options.BuildSettings = api.VersionedSettingsBuildSettings(d.Get("build_settings").(string))
	}
	if d.HasChange("context_parameters") {
		contextParametersRaw := d.Get("context_parameters").(map[string]interface{})
		vcsFeature.Options.ContextParameters = expandContextParameters(contextParametersRaw)
	}
	if d.HasChange("credentials_storage_type") {
		v := d.Get("credentials_storage_type").(string)
		if v == string(api.CredentialsStorageTypeCredentialsJSON) {
			vcsFeature.Options.CredentialsStorageType = api.CredentialsStorageTypeCredentialsJSON
		} else {
			vcsFeature.Options.CredentialsStorageType = api.CredentialsStorageTypeScrambledInVcs
		}
	}
	if d.HasChange("enabled") {
		vcsFeature.Options.Enabled = d.Get("enabled").(bool)
	}
	if d.HasChange("format") {
		vcsFeature.Options.Format = api.VersionedSettingsFormat(d.Get("format").(string))
	}
	if d.HasChange("show_changes") {
		vcsFeature.Options.ShowChanges = d.Get("show_changes").(bool)
	}
	if d.HasChange("use_relative_ids") {
		vcsFeature.Options.UseRelativeIds = d.Get("use_relative_ids").(bool)
	}
	if d.HasChange("vcs_root_id") {
		vcsFeature.Options.VcsRootID = d.Get("vcs_root_id").(string)
	}

	if _, err := service.Update(vcsFeature); err != nil {
		return err
	}

	return resourceProjectFeatureSlackNotifierRead(d, meta)
}

func resourceProjectFeatureSlackNotifierRead(d *schema.ResourceData, meta interface{}) error {
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

	slackNotifier, ok := feature.(*api.ProjectFeatureSlackNotifier)
	if !ok {
		return fmt.Errorf("Expected a VersionedSettings Feature but wasn't!")
	}

	d.Set("client_id", string(slackNotifier.Options.ClientId))
	d.Set("client_secret", d.Get("client_secret"))
	d.Set("display_name", string(slackNotifier.Options.DisplayName))
	d.Set("token", d.Get("token"))

	return nil
}

func resourceProjectFeatureSlackNotifierDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	projectId := d.Get("project_id").(string)
	featureId := d.Id()

	service := client.ProjectFeatureService(projectId)
	feature, err := service.GetByID(featureId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			// already gone
			return nil
		}

		return err
	}

	return service.Delete(feature.ID())
}

func resourceProjectFeatureSlackNotifierImportfunc(context.Context, *schema.ResourceData, interface{}) ([]*schema.ResourceData, error)
