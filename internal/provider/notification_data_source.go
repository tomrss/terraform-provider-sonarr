package provider

import (
	"context"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const notificationDataSourceName = "notification"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NotificationDataSource{}

func NewNotificationDataSource() datasource.DataSource {
	return &NotificationDataSource{}
}

// NotificationDataSource defines the notification implementation.
type NotificationDataSource struct {
	client *sonarr.APIClient
	auth   context.Context
}

func (d *NotificationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationDataSourceName
}

func (d *NotificationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Notifications -->\nSingle [Notification](../resources/notification).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Computed:            true,
			},
			"on_download": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Computed:            true,
			},
			"on_import_complete": schema.BoolAttribute{
				MarkdownDescription: "On import complete flag.",
				Computed:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Computed:            true,
			},
			"on_rename": schema.BoolAttribute{
				MarkdownDescription: "On rename flag.",
				Computed:            true,
			},
			"on_series_add": schema.BoolAttribute{
				MarkdownDescription: "On series add flag.",
				Computed:            true,
			},
			"on_series_delete": schema.BoolAttribute{
				MarkdownDescription: "On series delete flag.",
				Computed:            true,
			},
			"on_episode_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete flag.",
				Computed:            true,
			},
			"on_episode_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete for upgrade flag.",
				Computed:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Computed:            true,
			},
			"on_health_restored": schema.BoolAttribute{
				MarkdownDescription: "On health restored flag.",
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Computed:            true,
			},
			"on_manual_interaction_required": schema.BoolAttribute{
				MarkdownDescription: "On manual interaction required flag.",
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "Notification configuration template.",
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Notification implementation name.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Notification name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
			},
			// Field values
			"always_update": schema.BoolAttribute{
				MarkdownDescription: "Always update flag.",
				Computed:            true,
			},
			"clean_library": schema.BoolAttribute{
				MarkdownDescription: "Clean library flag.",
				Computed:            true,
			},
			"direct_message": schema.BoolAttribute{
				MarkdownDescription: "Direct message flag.",
				Computed:            true,
			},
			"notify": schema.BoolAttribute{
				MarkdownDescription: "Notify flag.",
				Computed:            true,
			},
			"send_silently": schema.BoolAttribute{
				MarkdownDescription: "Add silently flag.",
				Computed:            true,
			},
			"update_library": schema.BoolAttribute{
				MarkdownDescription: "Update library flag.",
				Computed:            true,
			},
			"use_eu_endpoint": schema.BoolAttribute{
				MarkdownDescription: "Use EU endpoint flag.",
				Computed:            true,
			},
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Computed:            true,
			},
			"use_encryption": schema.Int64Attribute{
				MarkdownDescription: "Require encryption.",
				Computed:            true,
			},
			"expire": schema.Int64Attribute{
				MarkdownDescription: "Expire.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Computed:            true,
			},
			"method": schema.Int64Attribute{
				MarkdownDescription: "Method. `1` POST, `2` PUT.",
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "Priority.",
				Computed:            true,
			},
			"retry": schema.Int64Attribute{
				MarkdownDescription: "Retry.",
				Computed:            true,
			},
			"notification_type": schema.Int64Attribute{
				MarkdownDescription: "Notification type. `0` Info, `1` Success, `2` Warning, `3` Failure.",
				Computed:            true,
			},
			"stateless_urls": schema.StringAttribute{
				MarkdownDescription: "Stateless URLs.",
				Computed:            true,
			},
			"configuration_key": schema.StringAttribute{
				MarkdownDescription: "Configuration key.",
				Computed:            true,
				Sensitive:           true,
			},
			"auth_username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Computed:            true,
				Sensitive:           true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Computed:            true,
			},
			"access_token_secret": schema.StringAttribute{
				MarkdownDescription: "Access token secret.",
				Computed:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Computed:            true,
				Sensitive:           true,
			},
			"app_token": schema.StringAttribute{
				MarkdownDescription: "App token.",
				Computed:            true,
				Sensitive:           true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Key.",
				Computed:            true,
				Sensitive:           true,
			},
			"event": schema.StringAttribute{
				MarkdownDescription: "Event.",
				Computed:            true,
			},
			"arguments": schema.StringAttribute{
				MarkdownDescription: "Arguments.",
				Computed:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "Author.",
				Computed:            true,
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Server URL.",
				Computed:            true,
			},
			"click_url": schema.StringAttribute{
				MarkdownDescription: "Click URL.",
				Computed:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Auth token.",
				Computed:            true,
				Sensitive:           true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Computed:            true,
			},
			"avatar": schema.StringAttribute{
				MarkdownDescription: "Avatar.",
				Computed:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
				Computed:            true,
			},
			"channel": schema.StringAttribute{
				MarkdownDescription: "Channel.",
				Computed:            true,
			},
			"chat_id": schema.StringAttribute{
				MarkdownDescription: "Chat ID.",
				Computed:            true,
			},
			"consumer_key": schema.StringAttribute{
				MarkdownDescription: "Consumer key.",
				Computed:            true,
			},
			"consumer_secret": schema.StringAttribute{
				MarkdownDescription: "Consumer secret.",
				Computed:            true,
				Sensitive:           true,
			},
			"device_names": schema.StringAttribute{
				MarkdownDescription: "Device names.",
				Computed:            true,
			},
			"display_time": schema.Int64Attribute{
				MarkdownDescription: "Display time.",
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Computed:            true,
			},
			"from": schema.StringAttribute{
				MarkdownDescription: "From.",
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Computed:            true,
			},
			"icon": schema.StringAttribute{
				MarkdownDescription: "Icon.",
				Computed:            true,
			},
			"mention": schema.StringAttribute{
				MarkdownDescription: "Mention.",
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "password.",
				Computed:            true,
				Sensitive:           true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path.",
				Computed:            true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Computed:            true,
			},
			"sender_domain": schema.StringAttribute{
				MarkdownDescription: "Sender domain.",
				Computed:            true,
			},
			"sender_id": schema.StringAttribute{
				MarkdownDescription: "Sender ID.",
				Computed:            true,
			},
			"sender_number": schema.StringAttribute{
				MarkdownDescription: "Sender Number.",
				Computed:            true,
				Sensitive:           true,
			},
			"receiver_id": schema.StringAttribute{
				MarkdownDescription: "Receiver ID.",
				Computed:            true,
			},
			"server": schema.StringAttribute{
				MarkdownDescription: "server.",
				Computed:            true,
			},
			"sign_in": schema.StringAttribute{
				MarkdownDescription: "Sign in.",
				Computed:            true,
			},
			"sound": schema.StringAttribute{
				MarkdownDescription: "Sound.",
				Computed:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Token.",
				Computed:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Computed:            true,
			},
			"user_key": schema.StringAttribute{
				MarkdownDescription: "User key.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"web_hook_url": schema.StringAttribute{
				MarkdownDescription: "Web hook url.",
				Computed:            true,
			},
			"grab_fields": schema.SetAttribute{
				MarkdownDescription: "Grab fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Group, `5` Size, `6` Links, `7` Release, `8` Poster, `9` Fanart.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"import_fields": schema.SetAttribute{
				MarkdownDescription: "Import fields. `0` Overview, `1` Rating, `2` Genres, `3` Quality, `4` Codecs, `5` Group, `6` Size, `7` Languages, `8` Subtitles, `9` Links, `10` Release, `11` Poster, `12` Fanart.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"channel_tags": schema.SetAttribute{
				MarkdownDescription: "Channel tags.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"device_ids": schema.SetAttribute{
				MarkdownDescription: "Device IDs.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: "Devices.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"recipients": schema.SetAttribute{
				MarkdownDescription: "Recipients.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"to": schema.SetAttribute{
				MarkdownDescription: "To.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"cc": schema.SetAttribute{
				MarkdownDescription: "Cc.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"bcc": schema.SetAttribute{
				MarkdownDescription: "Bcc.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"topics": schema.SetAttribute{
				MarkdownDescription: "Topics.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"field_tags": schema.SetAttribute{
				MarkdownDescription: "Tags and emojis.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *NotificationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *NotificationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *Notification

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get notification current value
	response, _, err := d.client.NotificationAPI.ListNotification(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationDataSourceName, err))

		return
	}

	data.find(ctx, data.Name.ValueString(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+notificationDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (n *Notification) find(ctx context.Context, name string, notifications []sonarr.NotificationResource, diags *diag.Diagnostics) {
	for _, notification := range notifications {
		if notification.GetName() == name {
			n.write(ctx, &notification, diags)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(notificationDataSourceName, "name", name))
}
