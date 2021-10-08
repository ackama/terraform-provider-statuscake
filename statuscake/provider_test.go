package statuscake_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	provider "terraform-provider-statuscake/statuscake"
	"testing"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"statuscake": func() (*schema.Provider, error) {
		return provider.New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := provider.New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
