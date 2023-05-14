package teamcity

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the plugin entry point
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"teamcity_artifact_dependency":                resourceArtifactDependency(),
			"teamcity_agent_requirement":                  resourceAgentRequirement(),
			"teamcity_build_config":                       resourceBuildConfig(),
			"teamcity_build_trigger_build_finish":         resourceBuildTriggerBuildFinish(),
			"teamcity_build_trigger_schedule":             resourceBuildTriggerSchedule(),
			"teamcity_build_trigger_vcs":                  resourceBuildTriggerVcs(),
			"teamcity_feature_commit_status_publisher":    resourceFeatureCommitStatusPublisher(),
			"teamcity_group":                              resourceGroup(),
			"teamcity_project":                            resourceProject(),
			"teamcity_project_feature_slack_notifier":     resourceProjectFeatureSlackNotifier(),
			"teamcity_project_feature_versioned_settings": resourceProjectFeatureVersionedSettings(),
			"teamcity_snapshot_dependency":                resourceSnapshotDependency(),
			"teamcity_vcs_root_git":                       resourceVcsRootGit(),
			"teamcity_agent_pool":                         resourceAgentPool(),
			"teamcity_agent_pool_project_assignment":      resourceAgentPoolProjectAssignment(),
			"teamcity_feature_golang":                     resourceFeatureGolang(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"teamcity_agent_pool": dataSourceAgentPool(),
			"teamcity_project":    dataSourceProject(),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TEAMCITY_ADDR", nil),
			},
			"token": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"username", "password"},
				DefaultFunc:   schema.EnvDefaultFunc("TEAMCITY_TOKEN", nil),
			},
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"token"},
				DefaultFunc:   schema.EnvDefaultFunc("TEAMCITY_USER", nil),
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"token"},
				DefaultFunc:   schema.EnvDefaultFunc("TEAMCITY_PASSWORD", nil),
			},
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Address:  d.Get("address").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	if v, ok := d.GetOk("token"); ok && v.(string) != "" {
		config.Token = v.(string)
	} else {
		config.Username = d.Get("username").(string)
		config.Password = d.Get("password").(string)
	}

	if config.Token == "" && config.Username == "" {
		return nil, fmt.Errorf("Error configuring provider: either a `token` or `username` must be specified")
	}

	return config.Client()
}
