---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_indexer_config Resource - terraform-provider-sonarr"
subcategory: ""
description: |-
  Indexer Config resource.For more information refer to Indexer https://wiki.servarr.com/sonarr/settings#options documentation.
---

# sonarr_indexer_config (Resource)

Indexer Config resource.<br/>For more information refer to [Indexer](https://wiki.servarr.com/sonarr/settings#options) documentation.

## Example Usage

```terraform
resource "sonarr_indexer_config" "example" {
  maximum_size      = 0
  minimum_age       = 0
  retention         = 0
  rss_sync_interval = 25
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `maximum_size` (Number) Maximum size.
- `minimum_age` (Number) Minimum age.
- `retention` (Number) Retention.
- `rss_sync_interval` (Number) RSS sync interval.

### Read-Only

- `id` (Number) Indexer Config ID.

## Import

Import is supported using the following syntax:

```shell
# import does not need parameters
terraform import sonarr_indexer_config.example
```