package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/sonarr"
)

const (
	downloadClientPneumaticResourceName   = "download_client_pneumatic"
	DownloadClientPneumaticImplementation = "Pneumatic"
	DownloadClientPneumaticConfigContrat  = "PneumaticSettings"
	DownloadClientPneumaticProtocol       = "usenet"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DownloadClientPneumaticResource{}
var _ resource.ResourceWithImportState = &DownloadClientPneumaticResource{}

func NewDownloadClientPneumaticResource() resource.Resource {
	return &DownloadClientPneumaticResource{}
}

// DownloadClientPneumaticResource defines the download client implementation.
type DownloadClientPneumaticResource struct {
	client *sonarr.Sonarr
}

// DownloadClientPneumatic describes the download client data model.
type DownloadClientPneumatic struct {
	Tags                     types.Set    `tfsdk:"tags"`
	Name                     types.String `tfsdk:"name"`
	NzbFolder                types.String `tfsdk:"nzb_folder"`
	StrmFolder               types.String `tfsdk:"strm_folder"`
	Priority                 types.Int64  `tfsdk:"priority"`
	ID                       types.Int64  `tfsdk:"id"`
	Enable                   types.Bool   `tfsdk:"enable"`
	RemoveFailedDownloads    types.Bool   `tfsdk:"remove_failed_downloads"`
	RemoveCompletedDownloads types.Bool   `tfsdk:"remove_completed_downloads"`
}

func (d DownloadClientPneumatic) toDownloadClient() *DownloadClient {
	return &DownloadClient{
		Tags:                     d.Tags,
		Name:                     d.Name,
		NzbFolder:                d.NzbFolder,
		StrmFolder:               d.StrmFolder,
		Priority:                 d.Priority,
		ID:                       d.ID,
		Enable:                   d.Enable,
		RemoveFailedDownloads:    d.RemoveFailedDownloads,
		RemoveCompletedDownloads: d.RemoveCompletedDownloads,
	}
}

func (d *DownloadClientPneumatic) fromDownloadClient(client *DownloadClient) {
	d.Tags = client.Tags
	d.Name = client.Name
	d.NzbFolder = client.NzbFolder
	d.StrmFolder = client.StrmFolder
	d.Priority = client.Priority
	d.ID = client.ID
	d.Enable = client.Enable
	d.RemoveFailedDownloads = client.RemoveFailedDownloads
	d.RemoveCompletedDownloads = client.RemoveCompletedDownloads
}

func (r *DownloadClientPneumaticResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + downloadClientPneumaticResourceName
}

func (r *DownloadClientPneumaticResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "<!-- subcategory:Download Clients -->Download Client Pneumatic resource.\nFor more information refer to [Download Client](https://wiki.servarr.com/sonarr/settings#download-clients) and [Pneumatic](https://wiki.servarr.com/sonarr/supported#pneumatic).",
		Attributes: map[string]tfsdk.Attribute{
			"enable": {
				MarkdownDescription: "Enable flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"remove_completed_downloads": {
				MarkdownDescription: "Remove completed downloads flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"remove_failed_downloads": {
				MarkdownDescription: "Remove failed downloads flag.",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			"priority": {
				MarkdownDescription: "Priority.",
				Optional:            true,
				Computed:            true,
				Type:                types.Int64Type,
			},
			"name": {
				MarkdownDescription: "Download Client name.",
				Required:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Optional:            true,
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Download Client ID.",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			// Field values
			"nzb_folder": {
				MarkdownDescription: "NZB folder.",
				Required:            true,
				Type:                types.StringType,
			},
			"strm_folder": {
				MarkdownDescription: "STRM folder.",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *DownloadClientPneumaticResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DownloadClientPneumaticResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var client *DownloadClientPneumatic

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new DownloadClientPneumatic
	request := client.read(ctx)

	response, err := r.client.AddDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to create %s, got error: %s", downloadClientPneumaticResourceName, err))

		return
	}

	tflog.Trace(ctx, "created "+downloadClientPneumaticResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientPneumaticResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var client DownloadClientPneumatic

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get DownloadClientPneumatic current value
	response, err := r.client.GetDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientPneumaticResourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+downloadClientPneumaticResourceName+": "+strconv.Itoa(int(response.ID)))
	// Map response body to resource schema attribute
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientPneumaticResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var client *DownloadClientPneumatic

	resp.Diagnostics.Append(req.Plan.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update DownloadClientPneumatic
	request := client.read(ctx)

	response, err := r.client.UpdateDownloadClientContext(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to update %s, got error: %s", downloadClientPneumaticResourceName, err))

		return
	}

	tflog.Trace(ctx, "updated "+downloadClientPneumaticResourceName+": "+strconv.Itoa(int(response.ID)))
	// Generate resource state struct
	client.write(ctx, response)
	resp.Diagnostics.Append(resp.State.Set(ctx, &client)...)
}

func (r *DownloadClientPneumaticResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var client *DownloadClientPneumatic

	resp.Diagnostics.Append(req.State.Get(ctx, &client)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete DownloadClientPneumatic current value
	err := r.client.DeleteDownloadClientContext(ctx, client.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", downloadClientPneumaticResourceName, err))

		return
	}

	tflog.Trace(ctx, "deleted "+downloadClientPneumaticResourceName+": "+strconv.Itoa(int(client.ID.ValueInt64())))
	resp.State.RemoveResource(ctx)
}

func (r *DownloadClientPneumaticResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			tools.UnexpectedImportIdentifier,
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)

		return
	}

	tflog.Trace(ctx, "imported "+downloadClientPneumaticResourceName+": "+strconv.Itoa(id))
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (d *DownloadClientPneumatic) write(ctx context.Context, downloadClient *sonarr.DownloadClientOutput) {
	genericDownloadClient := DownloadClient{
		Enable:                   types.BoolValue(downloadClient.Enable),
		RemoveCompletedDownloads: types.BoolValue(downloadClient.RemoveCompletedDownloads),
		RemoveFailedDownloads:    types.BoolValue(downloadClient.RemoveFailedDownloads),
		Priority:                 types.Int64Value(int64(downloadClient.Priority)),
		ID:                       types.Int64Value(downloadClient.ID),
		Name:                     types.StringValue(downloadClient.Name),
		Tags:                     types.SetValueMust(types.Int64Type, nil),
	}
	tfsdk.ValueFrom(ctx, downloadClient.Tags, genericDownloadClient.Tags.Type(ctx), &genericDownloadClient.Tags)
	genericDownloadClient.writeFields(ctx, downloadClient.Fields)
	d.fromDownloadClient(&genericDownloadClient)
}

func (d *DownloadClientPneumatic) read(ctx context.Context) *sonarr.DownloadClientInput {
	var tags []int

	tfsdk.ValueAs(ctx, d.Tags, &tags)

	return &sonarr.DownloadClientInput{
		Enable:                   d.Enable.ValueBool(),
		RemoveCompletedDownloads: d.RemoveCompletedDownloads.ValueBool(),
		RemoveFailedDownloads:    d.RemoveFailedDownloads.ValueBool(),
		Priority:                 int(d.Priority.ValueInt64()),
		ID:                       d.ID.ValueInt64(),
		ConfigContract:           DownloadClientPneumaticConfigContrat,
		Implementation:           DownloadClientPneumaticImplementation,
		Name:                     d.Name.ValueString(),
		Protocol:                 DownloadClientPneumaticProtocol,
		Tags:                     tags,
		Fields:                   d.toDownloadClient().readFields(ctx),
	}
}
