# terraform-provider-statuscake

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using
Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform
provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```hcl
terraform {
  required_providers {
    statuscake = {
      version = "0.1"
      source  = "ackama.com/ackama/statuscake"
    }
  }
}

provider "statuscake" {}

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

resource "statuscake_uptime_test" "my_site" {
  name        = "My Site"
  website_url = "https://www.example.com"
  test_type   = "HTTP"
  check_rate  = 300
  tags        = ["env:production", "app:example"]

  contact_groups = [
    statuscake_contact_group.main_contacts.id
  ]
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need
[Go](http://www.golang.org) installed on your machine (see
[Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put
the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
