package statuscake_test

import (
	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	provider "terraform-provider-statuscake/statuscake"
	"testing"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"statuscake": func() (*schema.Provider, error) { //nolint:unparam
		return provider.New("dev")(), nil
	},
}

// testAccPreCheck validates the necessary test API keys exist in the testing environment
func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("STATUSCAKE_API_KEY"); v == "" {
		t.Fatal("STATUSCAKE_API_KEY must be set for acceptance tests")
	}
}

func statusCakeAPIClient() *statuscake.APIClient {
	apiToken := os.Getenv("STATUSCAKE_API_KEY")

	return statuscake.NewAPIClient(apiToken)
}

func TestProvider(t *testing.T) {
	t.Parallel()

	if err := provider.New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
