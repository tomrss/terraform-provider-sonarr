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
	notificationTwitterResourceName   = "notification_twitter"
	NotificationTwitterImplementation = "Twitter"
	NotificationTwitterConfigContrat  = "TwitterSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationTwitterResource{}
var _ resource.ResourceWithImportState = &NotificationTwitterResource{}

func NewNotificationTwitterResource() resource.Resource {
	return &NotificationTwitterResource{}
}

// NotificationTwitterResource defines the notification implementation.
type NotificationTwitterResource struct {
	client *sonarr.Sonarr
}

// NotificationTwitter describes the notification data model.
type NotificationTwitter struct {
	Tags                          types.Set    `tfsdk:"tags"`
	Name                          types.String `tfsdk:"name"`
	AccessToken                   types.String `tfsdk:"access_token"`
	AccessTokenSecret             types.String `tfsdk:"access_token_secret"`
	ConsumerKey                   types.String `tfsdk:"consumer_key"`
	ConsumerSecret                types.String `tfsdk:"consumer_secret"`
	Mention                       types.String `tfsdk:"mention"`
	ID                            types.Int64  `tfsdk:"id"`
	DirectMessage                 types.Bool   `tfsdk:"direct_message"`
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

func (n NotificationTwitter) toNotification() *Notification {
	return &Notification{
		Tags:                          n.Tags,
		AccessToken:                   n.AccessToken,
		AccessTokenSecret:             n.AccessTokenSecret,
		ConsumerKey:                   n.ConsumerKey,
		ConsumerSecret:                n.ConsumerSecret,
		Mention:                       n.Mention,
		Name:                          n.Name,
		ID:                            n.ID,
		DirectMessage:                 n.DirectMessage,
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

func (n *NotificationTwitter) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.AccessToken = notification.AccessToken
	n.AccessTokenSecret = notification.AccessTokenSecret
	n.ConsumerKey = notification.ConsumerKey
	n.ConsumerSecret = notification.ConsumerSecret
	n.Mention = notification.Mention
	n.Name = notification.Name
	n.ID = notification.ID
	n.DirectMessage = notification.DirectMessage
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

func (r *NotificationTwitterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationTwitterResourceName
}

func (r *NotificationTwitterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->Notification Twitter resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Twitter](https://wiki.servarr.com/sonarr/supported#twitter).",
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
				MarkdownDescription: "NotificationTwitter name.",
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
			"direct_message": schema.BoolAttribute{
				MarkdownDescription: "Direct message flag.",
				Optional:            true,
				Computed:            true,
			},
			"consumer_key": schema.StringAttribute{
				MarkdownDescription: "Consumer Key.",
				Required:            true,
				Sensitive:           true,
			},
			"consumer_secret": schema.StringAttribute{
				MarkdownDescription: "Consumer Secret.",
				Required:            true,
				Sensitive:           true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Required:            true,
				Sensitive:           true,
			},
			"access_token_secret": schema.StringAttribute{
				MarkdownDescription: "Access token secret.",
				Required:            true,
				Sensitive:           true,
			},
			"mention": schema.StringAttribute{
				MarkdownDescription: "Mention.",
				Required:            true,
			},
		},
	}
}

func (r *NotificationTwitterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationTwitterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationTwitter

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationTwitter
	request := notification.read(ctx)

	response, err := r.client.AddNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", notificationTwitterResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationTwitterResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTwitterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationTwitter

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationTwitter current value
	response, err := r.client.GetNotificationContext(ctx, int(notification.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationTwitterResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationTwitterResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTwitterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationTwitter

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationTwitter
	request := notification.read(ctx)

	response, err := r.client.UpdateNotificationContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", notificationTwitterResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationTwitterResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	notification.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationTwitterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var notification *NotificationTwitter

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationTwitter current value
	err := r.client.DeleteNotificationContext(ctx, notification.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", notificationTwitterResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationTwitterResourceName+": "+strconv.Itoa(int(notification.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationTwitterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+notificationTwitterResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (n *NotificationTwitter) write(ctx context.Context, notification *sonarr.NotificationOutput) {
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

func (n *NotificationTwitter) read(ctx context.Context) *sonarr.NotificationInput {
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
		ConfigContract:                NotificationTwitterConfigContrat,
		Implementation:                NotificationTwitterImplementation,
		ID:                            n.ID.ValueInt64(),
		Name:                          n.Name.ValueString(),
		Tags:                          tags,
		Fields:                        n.toNotification().readFields(ctx),
	}
}