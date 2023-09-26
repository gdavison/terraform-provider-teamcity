package teamcity_test

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/cvbarros/go-teamcity/teamcity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccProjectFeatureSlackConnection_Basic(t *testing.T) {
	resourceName := "teamcity_project_feature_slack_notifier.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectFeatureSlackConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectFeatureSlackConnectionBasicConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProjectFeatureSlackConnectionExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "project_id", "teamcity_project.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "abcd.1234"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "xyz"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Notifier"),
					resource.TestCheckResourceAttr(resourceName, "token", "ABCD1234EFG"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDeploymentImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"token",
				},
			},
		},
	})
}

func TestAccProjectFeatureSlackConnection_Update(t *testing.T) {
	resourceName := "teamcity_project_feature_slack_notifier.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckProjectFeatureSlackConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectFeatureSlackConnectionBasicConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProjectFeatureSlackConnectionExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "project_id", "teamcity_project.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "abcd.1234"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "xyz"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Notifier"),
					resource.TestCheckResourceAttr(resourceName, "token", "ABCD1234EFG"),
				),
			},
			{
				Config: testAccProjectFeatureSlackConnectionUpdatedConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProjectFeatureSlackConnectionExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "project_id", "teamcity_project.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "client_id", "1234.abcd"),
					resource.TestCheckResourceAttr(resourceName, "client_secret", "abc"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated"),
					resource.TestCheckResourceAttr(resourceName, "token", "XYZ789ABCD"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDeploymentImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"client_secret",
					"token",
				},
			},
		},
	})
}

func testAccDeploymentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", resourceName)
		}

		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["project_id"], rs.Primary.ID), nil
	}
}

func testAccCheckProjectFeatureSlackConnectionExists(resourceName string) resource.TestCheckFunc {
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

		if _, ok := feature.(*api.ProjectFeatureSlackConnection); !ok {
			return fmt.Errorf("Expected a Versioned Setting but it wasn't!")
		}

		return nil
	}
}

func testAccCheckProjectFeatureSlackConnectionDestroy(s *terraform.State) error {
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

func testAccProjectFeatureSlackConnectionBasicConfig() string {
	return composeConfig(
		testAccProjectFeatureSlackConnectionTemplate, `
resource "teamcity_project_feature_slack_notifier" "test" {
  project_id    = teamcity_project.test.id
  client_id     = "abcd.1234"
  client_secret = "xyz"
  display_name  = "Notifier"
  token         = "ABCD1234EFG"
}
`)
}

func testAccProjectFeatureSlackConnectionUpdatedConfig() string {
	return composeConfig(
		testAccProjectFeatureSlackConnectionTemplate, `
resource "teamcity_project_feature_slack_notifier" "test" {
  project_id    = teamcity_project.test.id
  client_id     = "1234.abcd"
  client_secret = "abc"
  display_name  = "Updated"
  token         = "XYZ789ABCD"
}
`)
}

const testAccProjectFeatureSlackConnectionTemplate = `
resource "teamcity_project" "test" {
  name = "Test Project"
}
`
