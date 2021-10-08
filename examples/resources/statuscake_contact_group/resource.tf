locals {
  # currently statuscake doesn't provide a public api for managing integrations,
  # so you must create & get their IDs from the admin panel
  slack_integration_id = "12345"
}

resource "statuscake_contact_group" "main_contacts" {
  name = "Main Contacts"

  email_addresses = [
    "humans@example.com"
  ]

  integrations = [
    local.slack_integration_id
  ]
}
