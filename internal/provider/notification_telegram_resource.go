package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	notificationTelegramResourceName   = "notification_telegram"
	notificationTelegramImplementation = "Telegram"
	notificationTelegramConfigContract = "TelegramSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationTelegramResource{}
	_ resource.ResourceWithImportState = &NotificationTelegramResource{}
)

func NewNotificationTelegramResource() resource.Resource {
	return &NotificationTelegramResource{}
}

// NotificationTelegramResource defines the notification implementation.
type NotificationTelegramResource struct {
	client *sonarr.Sonarr
}

// NotificationTelegram describes the notification data model.
type NotificationTelegram struct {
	Tags                          types.Set    `tfsdk:"tags"`
	ChatID                        types.String `tfsdk:"chat_id"`
	Name                          types.String `tfsdk:"name"`
	BotToken                      types.String `tfsdk:"bot_token"`
	ID                            types.Int64  `tfsdk:"id"`
	SendSilently                  types.Bool   `tfsdk:"send_silently"`
	OnGrab                        types.Bool   `tfsdk:"on_grab"`
	OnEpisodeFileDeleteForUpgrade types.Bool   `tfsdk:"on_episode_file_delete_for_upgrade"`
	OnEpisodeFileDelete           types.Bool   `tfsdk:"on_episode_file_delete"`
	IncludeHealthWarnings         types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate           types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue                 types.Bool   `tfsdk:"on_health_issue"`
	OnSeriesDelete                types.Bool   `tfsdk:"on_series_delete"`
	OnUpgrade                     types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                    types.Bool   `tfsdk:"on_download"`
}

func (n NotificationTelegram) toNotification() *Notification {
	return &Notification{
		Tags:                          n.Tags,
		ChatID:                        n.ChatID,
		BotToken:                      n.BotToken,
		SendSilently:                  n.SendSilently,
		Name:                          n.Name,
		ID:                            n.ID,
		OnGrab:                        n.OnGrab,
		OnEpisodeFileDeleteForUpgrade: n.OnEpisodeFileDeleteForUpgrade,
		OnEpisodeFileDelete:           n.OnEpisodeFileDelete,
		IncludeHealthWarnings:         n.IncludeHealthWarnings,
		OnApplicationUpdate:           n.OnApplicationUpdate,
		OnHealthIssue:                 n.OnHealthIssue,
		OnSeriesDelete:                n.OnSeriesDelete,
		OnUpgrade:                     n.OnUpgrade,
		OnDownload:                    n.OnDownload,
	}
}

func (n *NotificationTelegram) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.ChatID = notification.ChatID
	n.BotToken = notification.BotToken
	n.SendSilently = notification.SendSilently
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnEpisodeFileDeleteForUpgrade = notification.OnEpisodeFileDeleteForUpgrade
	n.OnEpisodeFileDelete = notification.OnEpisodeFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnSeriesDelete = notification.OnSeriesDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
}

func (r *NotificationTelegramResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationTelegramResourceName
}

func (r *NotificationTelegramResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Telegram resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Telegram](https://wiki.servarr.com/sonarr/supported#telegram).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Required:            true,
			},
			"on_download": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Required:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Required:            true,
			},
			"on_series_delete": schema.BoolAttribute{
				MarkdownDescription: "On series delete flag.",
				Required:            true,
			},
			"on_episode_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete flag.",
				Required:            true,
			},
			"on_episode_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete for upgrade flag.",
				Required:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Required:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Required:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationTelegram name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Notification ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			// Field values
			"send_silently": schema.BoolAttribute{
				MarkdownDescription: "Send silently flag.",
				Optional:            true,
				Computed:            true,
			},
			"chat_id": schema.StringAttribute{
				MarkdownDescription: "Chat ID.",
				Required:            true,
			},
			"bot_token": schema.StringAttribute{
				MarkdownDescription: "Bot token.",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationTelegramResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedResourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NotificationTelegramResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationTelegram
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationTelegram current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationTelegram
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationTelegramResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTelegramResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationTelegram

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationTelegram current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationTelegramResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationTelegramResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationTelegramResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationTelegramResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationTelegram) write(ctx context.Context, notification *sonarr.NotificationOutput) {
	genericNotification := Notification{
		OnGrab:                        types.BoolValue(notification.OnGrab),
		OnDownload:                    types.BoolValue(notification.OnDownload),
		OnUpgrade:                     types.BoolValue(notification.OnUpgrade),
		OnSeriesDelete:                types.BoolValue(notification.OnSeriesDelete),
		OnEpisodeFileDelete:           types.BoolValue(notification.OnEpisodeFileDelete),
		OnEpisodeFileDeleteForUpgrade: types.BoolValue(notification.OnEpisodeFileDeleteForUpgrade),
		OnHealthIssue:                 types.BoolValue(notification.OnHealthIssue),
		OnApplicationUpdate:           types.BoolValue(notification.OnApplicationUpdate),
		IncludeHealthWarnings:         types.BoolValue(notification.IncludeHealthWarnings),
		ID:                            types.Int64Value(notification.ID),
		Name:                          types.StringValue(notification.Name),
	}
	genericNotification.Tags, _ = types.SetValueFrom(ctx, types.Int64Type, notification.Tags)
	genericNotification.writeFields(ctx, notification.Fields)
	n.fromNotification(&genericNotification)
}

func (n *NotificationTelegram) read(ctx context.Context) *sonarr.NotificationInput {
	var tags []int

	tfsdk.ValueAs(ctx, n.Tags, &tags)

	return &sonarr.NotificationInput{
		OnGrab:                        n.OnGrab.ValueBool(),
		OnDownload:                    n.OnDownload.ValueBool(),
		OnUpgrade:                     n.OnUpgrade.ValueBool(),
		OnSeriesDelete:                n.OnSeriesDelete.ValueBool(),
		OnEpisodeFileDelete:           n.OnEpisodeFileDelete.ValueBool(),
		OnEpisodeFileDeleteForUpgrade: n.OnEpisodeFileDeleteForUpgrade.ValueBool(),
		OnHealthIssue:                 n.OnHealthIssue.ValueBool(),
		OnApplicationUpdate:           n.OnApplicationUpdate.ValueBool(),
		IncludeHealthWarnings:         n.IncludeHealthWarnings.ValueBool(),
		ConfigContract:                notificationTelegramConfigContract,
		Implementation:                notificationTelegramImplementation,
		ID:                            n.ID.ValueInt64(),
		Name:                          n.Name.ValueString(),
		Tags:                          tags,
		Fields:                        n.toNotification().readFields(ctx),
	}
}
