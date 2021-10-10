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

func fetchAllUptimeTests() ([]statuscake.UptimeTestOverview, error) {
	var uptimeTests []statuscake.UptimeTestOverview

	client := statusCakeAPIClient()
	currentPage := int32(1)

	for {
		res, err := client.ListUptimeTests(context.TODO()).
			Page(currentPage).
			Execute()

		if err != nil {
			return nil, fmt.Errorf("failed to fetch uptime tests: %w", err)
		}

		uptimeTests = append(uptimeTests, res.Data...)

		currentPage += 1 //nolint:revive

		if currentPage >= *res.Metadata.PageCount {
			return uptimeTests, nil
		}
	}
}

// testAccCheckUptimeTestDestroy verifies the uptime test has been destroyed
func testAccCheckUptimeTestDestroy(s *terraform.State) error {
	// loop through the resources in state, verifying each uptime test is destroyed
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "statuscake_uptime_test" {
			continue
		}

		uptimeTests, err := fetchAllUptimeTests()

		if err == nil {
			if len(uptimeTests) > 0 {
				for _, uptimeTest := range uptimeTests {
					if uptimeTest.ID == rs.Primary.ID {
						return fmt.Errorf("uptime test (%s) still exists", rs.Primary.ID)
					}
				}
			}

			return nil
		}
	}

	return nil
}

func testAccCheckUptimeTestExists(resourceName string) resource.TestCheckFunc { //nolint:unparam
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("uptime test ID is not set")
		}

		// fetch *all* the uptime tests to be sure we're using a unique id
		uptimeTests, err := fetchAllUptimeTests()

		if err != nil {
			return err
		}

		finds := 0

		for _, uptimeTest := range uptimeTests {
			if uptimeTest.ID == rs.Primary.ID {
				finds += 1 //nolint:revive
			}
		}

		if finds == 0 {
			return fmt.Errorf("uptime test not found")
		}
		if finds >= 2 {
			return fmt.Errorf("multiple uptime tests matching id found")
		}

		return nil
	}
}

func TestAccUptimeTest_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckUptimeTestDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name        = "My Site"
						website_url = "https://www.example.com"
						test_type   = "HTTP"
            check_rate  = 300
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUptimeTestExists("statuscake_uptime_test.foo"),
				),
			},
		},
	})
}

func TestAccUptimeTest_changing(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckUptimeTestDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name             = "My Site"
						website_url      = "https://www.example.com"
						test_type        = "HTTP"
						check_rate       = 300
						confirmation     = 3
					  custom_header    = ""
						do_not_find      = true
						dns_server       = "my-host"
						enable_ssl_alert = true
						final_endpoint   = "https://www.example.com"
						find_string      = "example"
						follow_redirects = true
						host             = "The World"
						paused           = true
						port             = 443
						post_body        = "{}"
						post_raw         = ""
						timeout          = 60
						trigger_rate     = 5
						cookie_storage   = true
						user_agent       = "StatusCake"
						status_codes 		 = ["401", "405", "406"]
						tags             = ["one", "two", "three"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUptimeTestExists("statuscake_uptime_test.foo"),
				),
			},
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name             = "My Site!"
						website_url      = "https://www.example.com"
						test_type        = "HTTP"
						check_rate       = 900
						confirmation     = 1
					  custom_header    = ""
						do_not_find      = false
						dns_server       = "my-other-host"
						enable_ssl_alert = false
						final_endpoint   = "https://www.example.com.au"
						find_string      = "no-thanks"
						follow_redirects = false
						host             = "The Moon"
						paused           = false
						port             = 80
						post_body        = ""
						post_raw         = "{}"
						timeout          = 30
						trigger_rate     = 10
						cookie_storage   = false
						user_agent       = "StatusCake2"
						status_codes 		 = ["401", "301"]
						tags             = ["five", "six"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUptimeTestExists("statuscake_uptime_test.foo"),
				),
			},
		},
	})
}

func TestAccUptimeTest_validation(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		IsUnitTest:        true,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckUptimeTestDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name             = "My Site"
						website_url      = "https://www.example.com"
						test_type        = "F2P"
						check_rate       = 300
					}
				`,
				ExpectError: regexp.MustCompile("expected test_type to be one of"),
			},
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name             = "My Site"
						website_url      = "https://www.example.com"
						test_type        = "HTTP"
						check_rate       = 234
					}
				`,
				ExpectError: regexp.MustCompile("expected check_rate to be one of"),
			},
		},
	})
}

func TestAccUptimeTest_prettyError(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckUptimeTestDestroy,
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name             = "My Site"
						website_url      = "https://www.example.com"
						test_type        = "HTTP"
						check_rate       = 300
						confirmation     = 5
					}
				`,
				ExpectError: regexp.MustCompile("Confirmation must be no more than 3"),
			},
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name             = "My Site"
						website_url      = "https://www.example.com"
						test_type        = "HTTP"
						check_rate       = 300
						confirmation     = 3
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUptimeTestExists("statuscake_uptime_test.foo"),
				),
			},
			{
				Config: `
					resource "statuscake_uptime_test" "foo" {
						name             = "My Site"
						website_url      = "https://www.example.com"
						test_type        = "HTTP"
						check_rate       = 300
						confirmation     = 5
					}
				`,
				ExpectError: regexp.MustCompile("Confirmation must be no more than 3"),
			},
		},
	})
}
