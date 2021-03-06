---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "statuscake_contact_group Resource - terraform-provider-statuscake"
subcategory: ""
description: |-
  Manages a StatusCake Contact Group
---

# statuscake_contact_group (Resource)

Manages a StatusCake Contact Group

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Name of the contact group

### Optional

- **email_addresses** (List of String) List of email addresses
- **id** (String) The ID of this resource.
- **integrations** (List of String) List of integration IDs
- **mobile_numbers** (List of String) List of international format mobile phone numbers
- **ping_url** (String) URL or IP address of an endpoint to push uptime events. Currently this only supports HTTP GET endpoints


