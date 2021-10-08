resource "statuscake_uptime_test" "my_site" {
  name        = "My Site"
  website_url = "https://www.example.com"
  test_type   = "HTTP"
  check_rate  = 300
  tags        = ["env:production", "app:example"]
}
