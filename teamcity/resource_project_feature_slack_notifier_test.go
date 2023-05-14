package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProjectFeatureSlackNotifier_Basic(t *testing.T) {
	resName := "teamcity_project_feature_slack_notifier.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectFeatureSlackNotifierDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectFeatureSlackNotifierBasicConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProjectFeatureSlackNotifierExists(resName),
					resource.TestCheckResourceAttrPair(resName, "project_id", "teamcity_project.test", "id"),
					resource.TestCheckResourceAttr(resName, "client_id", "abcd.1234"),
					resource.TestCheckResourceAttr(resName, "client_secret", "xyz"),
					resource.TestCheckResourceAttr(resName, "display_name", "Notifier"),
					resource.TestCheckResourceAttr(resName, "token", "ABCD1234EFG"),
				),
			},
			// {
			// 	ResourceName:      resName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func testAccCheckProjectFeatureSlackNotifierExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*api.Client)
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		projectId := rs.Primary.Attributes["project_id"]
		featureId := rs.Primary.ID

		service := client.ProjectFeatureService(projectId)
		feature, err := service.GetByID(featureId)
		if err != nil {
			return fmt.Errorf("Received an error retrieving project versioned settings: %s", err)
		}

		if _, ok := feature.(*api.ProjectFeatureSlackNotifier); !ok {
			return fmt.Errorf("Expected a Versioned Setting but it wasn't!")
		}

		return nil
	}
}

func testAccCheckProjectFeatureSlackNotifierDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "teamcity_project_feature_slack_notifier" {
			continue
		}

		projectId := r.Primary.Attributes["project_id"]
		featureId := r.Primary.ID

		service := client.ProjectFeatureService(projectId)
		if _, err := service.GetByID(featureId); err != nil {
			if strings.Contains(err.Error(), "404") {
				// expected, since it's gone
				continue
			}

			return fmt.Errorf("Received an error retrieving project versioned settings: %s", err)
		}

		return fmt.Errorf("Project Versioned Settings still exists")
	}
	return nil
}

func testAccProjectFeatureSlackNotifierBasicConfig() string {
	return fmt.Sprintf(`
%[1]s

resource "teamcity_project_feature_slack_notifier" "test" {
  project_id     = teamcity_project.test.id
  client_id    = "abcd.1234"
  client_secret = "xyz"
  display_name         = "Notifier"
  token = "ABCD1234EFG"
}
`, testAccProjectFeatureSlackNotifierTemplate)
}

const testAccProjectFeatureSlackNotifierTemplate = `
resource "teamcity_project" "test" {
  name = "Test Project"
}
`
