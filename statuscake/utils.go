package statuscake

import (
	"encoding/json"
	"errors"
	"github.com/StatusCakeDev/statuscake-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
)

func asListOfStrings(list interface{}) []string {
	strings := make([]string, 0, len(list.([]interface{})))

	for _, item := range list.([]interface{}) {
		strings = append(strings, item.(string))
	}

	return strings
}

func apiErrorDiag(err error) diag.Diagnostics {
	var apiError statuscake.APIError
	var diags diag.Diagnostics

	if !errors.As(err, &apiError) {
		return diag.Errorf("Unknown error (received error that unexpectedly not an api error)")
	}

	for _, errs := range apiError.Errors {
		for _, message := range errs {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  message,
			})
		}
	}

	return diags
}

func isNotFoundAPIError(err error) bool {
	var apiError statuscake.APIError

	return errors.As(err, &apiError) && apiError.Status == 404
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
	var statuscakeError statuscake.APIError

	if !errors.As(err, &statuscakeError) {
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
