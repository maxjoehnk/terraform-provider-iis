# Terraform IIS Provider

Terraform Provider using the [Microsoft IIS Administration](https://docs.microsoft.com/en-us/IIS-Administration/) API.

## Usage
```hcl
provider "iis" {
  access_key = "your access key"
  host = "https://localhost:55539"
}

resource "iis_application_pool" "name" {
  name = "AppPool" // Name of the Application Pool
}

resource "iis_application" "name" {
  physical_path = "%systemdrive%\\inetpub\\your_app" // Path on the server to your web app
  application_pool = "${iis_application_pool.name.id}"
  path = "YourApp" // Path for URL access
  website = "${data.iis_website.default.ids[0]}" // id for the website is required
}

data "iis_website" "default" {}
```