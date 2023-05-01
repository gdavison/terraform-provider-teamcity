package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/cvbarros/terraform-provider-teamcity/teamcity"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: teamcity.Provider})
}
