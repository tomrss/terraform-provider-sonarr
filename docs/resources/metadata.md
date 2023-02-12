---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_metadata Resource - terraform-provider-sonarr"
subcategory: "Metadata"
description: |-
  Generic Metadata resource. When possible use a specific resource instead.
  For more information refer to Metadata https://wiki.servarr.com/sonarr/settings#metadata documentation.
---

# sonarr_metadata (Resource)

<!-- subcategory:Metadata -->Generic Metadata resource. When possible use a specific resource instead.
For more information refer to [Metadata](https://wiki.servarr.com/sonarr/settings#metadata) documentation.

## Example Usage

```terraform
resource "sonarr_metadata" "example" {
  enable           = true
  name             = "Example"
  implementation   = "MediaBrowserMetadata"
  config_contract  = "MediaBrowserMetadataSettings"
  episode_metadata = true
  series_images    = false
  season_images    = true
  episode_images   = false
  tags             = [1, 2]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config_contract` (String) Metadata configuration template.
- `implementation` (String) Metadata implementation name.
- `name` (String) Metadata name.

### Optional

- `enable` (Boolean) Enable flag.
- `episode_images` (Boolean) Episode images flag.
- `episode_metadata` (Boolean) Episode metadata flag.
- `season_images` (Boolean) Season images flag.
- `series_images` (Boolean) Series images flag.
- `series_metadata` (Boolean) Series metafata flag.
- `series_metadata_url` (Boolean) Series metadata URL flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Metadata ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_metadata.example 1
```