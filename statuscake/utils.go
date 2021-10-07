package statuscake

import (
	"encoding/json"
	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
)

func asListOfStrings(list interface{}) []string {
	var strings []string

	for _, item := range list.([]interface{}) {
		strings = append(strings, item.(string))
	}

	return strings
}

func asDiag(apiError statuscake.APIError) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, errors := range apiError.Errors {
		for _, message := range errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  message,
			})
		}
	}

	return diags
}

func logResponse(res interface{}) {
	log.Printf("[DEBUG] StatusCake API Response: %s", prettifyObject(res))
}

func prettifyObject(obj interface{}) string {
	pretty, err := json.MarshalIndent(obj, "", "  ")

	if err != nil {
		log.Printf("[WARN] Failed to make object pretty: %s", err)

		return "<ERROR>"
	}

	return string(pretty)
}

func logStatusCakeAPIError(err error) {
	statuscakeError, ok := err.(statuscake.APIError)

	if !ok {
		log.Printf("[WARN] Was provided an error that was not from StatusCake API")
		log.Printf("[DEBUG] Error: %s", err)

		return
	}

	pretty, err := json.MarshalIndent(statuscakeError, "", "  ")

	if err != nil {
		log.Printf("[WARN] Failed to make StatusCake API error pretty: %s", err)
	} else {
		log.Printf("[DEBUG] StatusCake API returned an error: %s", pretty)
	}
}
