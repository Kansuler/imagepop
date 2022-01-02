# terraform-provider-imagepop

![License](https://img.shields.io/github/license/Kansuler/imagepop) ![Version](https://img.shields.io/github/go-mod/go-version/Kansuler/imagepop)

This terraform provider will check that an docker image with a certain tag has been pushed to the docker registry. This can be useful if you want terraform to update a resource, but want terraform to wait until a separate build pipeline has finished.

### Usage

```hcl
terraform {
  required_providers {
    imagepop = {
      version = "0.1.0"
      source  = "Kansuler/imagepop"
    }
  }
}

data "imagepop" "example" {
  registry   = "https://gcr.io"
  repository = "Kansuler/example"
  tag        = <tag> // Tag to check for existance
  username   = <username>
  password   = <credentials>

  retry {
    attempts = 30 // How many attempts
    delay    = 3 // Delay between each attempt in seconds
  }
}
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
