package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/sonarr-go/sonarr"
	"github.com/devopsarr/terraform-provider-sonarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	notificationSignalResourceName   = "notification_signal"
	notificationSignalImplementation = "Signal"
	notificationSignalConfigContract = "SignalSettings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &NotificationSignalResource{}
	_ resource.ResourceWithImportState = &NotificationSignalResource{}
)

func NewNotificationSignalResource() resource.Resource {
	return &NotificationSignalResource{}
}

// NotificationSignalResource defines the notification implementation.
type NotificationSignalResource struct {
	client *sonarr.APIClient
	auth   context.Context
}

// NotificationSignal describes the notification data model.
type NotificationSignal struct {
	Tags                          types.Set    `tfsdk:"tags"`
	AuthPassword                  types.String `tfsdk:"auth_password"`
	AuthUsername                  types.String `tfsdk:"auth_username"`
	Host                          types.String `tfsdk:"host"`
	SenderNumber                  types.String `tfsdk:"sender_number"`
	ReceiverID                    types.String `tfsdk:"receiver_id"`
	Name                          types.String `tfsdk:"name"`
	Port                          types.Int64  `tfsdk:"port"`
	ID                            types.Int64  `tfsdk:"id"`
	UseSSL                        types.Bool   `tfsdk:"use_ssl"`
	OnGrab                        types.Bool   `tfsdk:"on_grab"`
	OnEpisodeFileDeleteForUpgrade types.Bool   `tfsdk:"on_episode_file_delete_for_upgrade"`
	OnEpisodeFileDelete           types.Bool   `tfsdk:"on_episode_file_delete"`
	IncludeHealthWarnings         types.Bool   `tfsdk:"include_health_warnings"`
	OnApplicationUpdate           types.Bool   `tfsdk:"on_application_update"`
	OnHealthIssue                 types.Bool   `tfsdk:"on_health_issue"`
	OnHealthRestored              types.Bool   `tfsdk:"on_health_restored"`
	OnManualInteractionRequired   types.Bool   `tfsdk:"on_manual_interaction_required"`
	OnSeriesAdd                   types.Bool   `tfsdk:"on_series_add"`
	OnSeriesDelete                types.Bool   `tfsdk:"on_series_delete"`
	OnUpgrade                     types.Bool   `tfsdk:"on_upgrade"`
	OnDownload                    types.Bool   `tfsdk:"on_download"`
	OnImportComplete              types.Bool   `tfsdk:"on_import_complete"`
}

func (n NotificationSignal) toNotification() *Notification {
	return &Notification{
		Tags:                          n.Tags,
		AuthPassword:                  n.AuthPassword,
		Name:                          n.Name,
		AuthUsername:                  n.AuthUsername,
		Host:                          n.Host,
		SenderNumber:                  n.SenderNumber,
		ReceiverID:                    n.ReceiverID,
		Port:                          n.Port,
		UseSSL:                        n.UseSSL,
		ID:                            n.ID,
		OnGrab:                        n.OnGrab,
		OnEpisodeFileDeleteForUpgrade: n.OnEpisodeFileDeleteForUpgrade,
		OnEpisodeFileDelete:           n.OnEpisodeFileDelete,
		IncludeHealthWarnings:         n.IncludeHealthWarnings,
		OnApplicationUpdate:           n.OnApplicationUpdate,
		OnHealthIssue:                 n.OnHealthIssue,
		OnHealthRestored:              n.OnHealthRestored,
		OnManualInteractionRequired:   n.OnManualInteractionRequired,
		OnSeriesAdd:                   n.OnSeriesAdd,
		OnSeriesDelete:                n.OnSeriesDelete,
		OnUpgrade:                     n.OnUpgrade,
		OnDownload:                    n.OnDownload,
		OnImportComplete:              n.OnImportComplete,
		ConfigContract:                types.StringValue(notificationSignalConfigContract),
		Implementation:                types.StringValue(notificationSignalImplementation),
	}
}

func (n *NotificationSignal) fromNotification(notification *Notification) {
	n.Tags = notification.Tags
	n.AuthPassword = notification.AuthPassword
	n.AuthUsername = notification.AuthUsername
	n.Host = notification.Host
	n.SenderNumber = notification.SenderNumber
	n.ReceiverID = notification.ReceiverID
	n.Port = notification.Port
	n.UseSSL = notification.UseSSL
	n.Name = notification.Name
	n.ID = notification.ID
	n.OnGrab = notification.OnGrab
	n.OnEpisodeFileDeleteForUpgrade = notification.OnEpisodeFileDeleteForUpgrade
	n.OnEpisodeFileDelete = notification.OnEpisodeFileDelete
	n.IncludeHealthWarnings = notification.IncludeHealthWarnings
	n.OnApplicationUpdate = notification.OnApplicationUpdate
	n.OnHealthIssue = notification.OnHealthIssue
	n.OnHealthRestored = notification.OnHealthRestored
	n.OnManualInteractionRequired = notification.OnManualInteractionRequired
	n.OnSeriesAdd = notification.OnSeriesAdd
	n.OnSeriesDelete = notification.OnSeriesDelete
	n.OnUpgrade = notification.OnUpgrade
	n.OnDownload = notification.OnDownload
	n.OnImportComplete = notification.OnImportComplete
}

func (r *NotificationSignalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + notificationSignalResourceName
}

func (r *NotificationSignalResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Notifications -->\nNotification Signal resource.\nFor more information refer to [Notification](https://wiki.servarr.com/sonarr/settings#connect) and [Signal](https://wiki.servarr.com/sonarr/supported#signal).",
		Attributes: map[string]schema.Attribute{
			"on_grab": schema.BoolAttribute{
				MarkdownDescription: "On grab flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_download": schema.BoolAttribute{
				MarkdownDescription: "On download flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_import_complete": schema.BoolAttribute{
				MarkdownDescription: "On import complete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_series_add": schema.BoolAttribute{
				MarkdownDescription: "On series add flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_series_delete": schema.BoolAttribute{
				MarkdownDescription: "On series delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_episode_file_delete": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_episode_file_delete_for_upgrade": schema.BoolAttribute{
				MarkdownDescription: "On episode file delete for upgrade flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_issue": schema.BoolAttribute{
				MarkdownDescription: "On health issue flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_health_restored": schema.BoolAttribute{
				MarkdownDescription: "On health restored flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_manual_interaction_required": schema.BoolAttribute{
				MarkdownDescription: "On manual interaction required flag.",
				Optional:            true,
				Computed:            true,
			},
			"on_application_update": schema.BoolAttribute{
				MarkdownDescription: "On application update flag.",
				Optional:            true,
				Computed:            true,
			},
			"include_health_warnings": schema.BoolAttribute{
				MarkdownDescription: "Include health warnings.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "NotificationSignal name.",
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
			"use_ssl": schema.BoolAttribute{
				MarkdownDescription: "Use SSL flag.",
				Optional:            true,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Optional:            true,
				Computed:            true,
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "Host.",
				Required:            true,
			},
			"sender_number": schema.StringAttribute{
				MarkdownDescription: "Sender Number.",
				Required:            true,
				Sensitive:           true,
			},
			"receiver_id": schema.StringAttribute{
				MarkdownDescription: "Receiver ID.",
				Required:            true,
			},
			"auth_username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Optional:            true,
				Computed:            true,
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: "Password.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *NotificationSignalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if auth, client := resourceConfigure(ctx, req, resp); client != nil {
		r.client = client
		r.auth = auth
	}
}

func (r *NotificationSignalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var notification *NotificationSignal

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new NotificationSignal
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.CreateNotification(r.auth).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Create, notificationSignalResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+notificationSignalResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSignalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var notification *NotificationSignal

	resp.Diagnostics.Append(req.State.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get NotificationSignal current value
	response, _, err := r.client.NotificationAPI.GetNotificationById(r.auth, int32(notification.ID.ValueInt64())).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, notificationSignalResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+notificationSignalResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Map response body to resource schema attribute
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSignalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var notification *NotificationSignal

	resp.Diagnostics.Append(req.Plan.Get(ctx, &notification)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update NotificationSignal
	request := notification.read(ctx, &resp.Diagnostics)

	response, _, err := r.client.NotificationAPI.UpdateNotification(r.auth, strconv.Itoa(int(request.GetId()))).NotificationResource(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Update, notificationSignalResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+notificationSignalResourceName+": "+strconv.Itoa(int(response.GetId())))
	// Generate resource state struct
	notification.write(ctx, response, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &notification)...)
}

func (r *NotificationSignalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var ID int64

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &ID)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete NotificationSignal current value
	_, err := r.client.NotificationAPI.DeleteNotification(r.auth, int32(ID)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Delete, notificationSignalResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+notificationSignalResourceName+": "+strconv.Itoa(int(ID)))
	resp.State.RemoveResource(ctx)
}

func (r *NotificationSignalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	helpers.ImportStatePassthroughIntID(ctx, path.Root("id"), req, resp)
	tflog.Trace(ctx, "imported "+notificationSignalResourceName+": "+req.ID)
}

func (n *NotificationSignal) write(ctx context.Context, notification *sonarr.NotificationResource, diags *diag.Diagnostics) {
	genericNotification := n.toNotification()
	genericNotification.write(ctx, notification, diags)
	n.fromNotification(genericNotification)
}

func (n *NotificationSignal) read(ctx context.Context, diags *diag.Diagnostics) *sonarr.NotificationResource {
	return n.toNotification().read(ctx, diags)
}
