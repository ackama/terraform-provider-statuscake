package statuscake_test

import (
	"context"
	"fmt"
	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

func fetchAllContactGroups() ([]statuscake.ContactGroup, error) {
	client := statusCakeAPIClient()

	res, err := client.ListContactGroups(context.TODO()).Execute()

	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

// testAccCheckContactGroupDestroy verifies the contact groups has been destroyed
func testAccCheckContactGroupDestroy(s *terraform.State) error {
	// loop through the resources in state, verifying each contact group is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "statuscake_contact_group" {
			continue
		}

		contactGroups, err := fetchAllContactGroups()

		if err == nil {
			if len(contactGroups) > 0 {
				for _, contactGroup := range contactGroups {
					if contactGroup.ID == rs.Primary.ID {
						return fmt.Errorf("contact group (%s) still exists", rs.Primary.ID)
					}
				}
			}

			return nil
		}
	}

	return nil
}

func testAccCheckContactGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("contact group ID is not set")
		}

		// fetch *all* the contact groups to be sure we're using a unique id
		contactGroups, err := fetchAllContactGroups()

		if err != nil {
			return err
		}

		finds := 0

		for _, contactGroup := range contactGroups {
			if contactGroup.ID == rs.Primary.ID {
				finds += 1
			}
		}

		if finds == 0 {
			return fmt.Errorf("contract group not found")
		}
		if finds >= 2 {
			return fmt.Errorf("multiple contact groups matching id found")
		}

		return nil
	}
}

func TestAccContactGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckContactGroupDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "statuscake_contact_group" "foo" {
						name = "My Group"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactGroupExists("statuscake_contact_group.foo"),
				),
			},
		},
	})
}

func TestAccContactGroup_prettyError(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckContactGroupDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "statuscake_contact_group" "foo" {
						name     = "My Group"
						ping_url = "not-a-url"
					}
				`,
				ExpectError: regexp.MustCompile("Ping Url is not a valid URL"),
			},
			{
				Config: `
					resource "statuscake_contact_group" "foo" {
						name     = "My Group"
						ping_url = "https://www.example.com"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactGroupExists("statuscake_contact_group.foo"),
				),
			},
			{
				Config: `
					resource "statuscake_contact_group" "foo" {
						name     = "My Group"
						ping_url = "not-a-url"
					}
				`,
				ExpectError: regexp.MustCompile("Ping Url is not a valid URL"),
			},
		},
	})
}
