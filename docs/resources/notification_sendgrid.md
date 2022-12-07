---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarr_notification_sendgrid Resource - terraform-provider-sonarr"
subcategory: "Notifications"
description: |-
  Notification Sendgrid resource.
  For more information refer to Notification https://wiki.servarr.com/sonarr/settings#connect and Sendgrid https://wiki.servarr.com/sonarr/supported#sendgrid.
---

# sonarr_notification_sendgrid (Resource)

<!-- subcategory:Notifications -->Notification Sendgrid resource.
For more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Sendgrid](https://wiki.servarr.com/sonarr/supported#sendgrid).

## Example Usage

```terraform
resource "sonarr_notification_sendgrid" "example" {
  on_grab                            = false
  on_download                        = true
  on_upgrade                         = true
  on_series_delete                   = false
  on_episode_file_delete             = false
  on_episode_file_delete_for_upgrade = true
  on_health_issue                    = false
  on_application_update              = false

  include_health_warnings = false
  name                    = "Example"

  api_key    = "APIkey"
  from       = "from_sendgrid@example.com"
  recipients = ["user1@example.com", "user2@example.com"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `from` (String) From.
- `include_health_warnings` (Boolean) Include health warnings.
- `name` (String) NotificationSendgrid name.
- `on_application_update` (Boolean) On application update flag.
- `on_download` (Boolean) On download flag.
- `on_episode_file_delete` (Boolean) On episode file delete flag.
- `on_episode_file_delete_for_upgrade` (Boolean) On episode file delete for upgrade flag.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_series_delete` (Boolean) On series delete flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `recipients` (Set of String) Recipients.

### Optional

- `api_key` (String, Sensitive) API key.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import sonarr_notification_sendgrid.example 1
```