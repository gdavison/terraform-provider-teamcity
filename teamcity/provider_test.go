package teamcity_test

import (
	"os"
	"testing"

	"github.com/cvbarros/terraform-provider-teamcity/teamcity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = teamcity.Provider()
	testAccProviders = map[string]*schema.Provider{
		"teamcity": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := teamcity.Provider().InternalValidate(); err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = teamcity.Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TEAMCITY_ADDR"); v == "" {
		t.Fatal("TEAMCITY_ADDR must be set for acceptance tests")
	}
	hasToken := os.Getenv("TEAMCITY_TOKEN") != ""
	hasUsername := os.Getenv("TEAMCITY_USER") != ""
	hasPassword := os.Getenv("TEAMCITY_PASSWORD") != ""

	if !hasToken && !(hasUsername && hasPassword) {
		t.Fatal("Either `TEAMCITY_TOKEN` or `TEAMCITY_USER` and `TEAMCITY_PASSWORD` must be set for acceptance tests")
	}
}
