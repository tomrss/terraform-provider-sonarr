package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr"
	"golift.io/starr/sonarr"
)

type resourceLanguageProfileType struct{}

type resourceLanguageProfile struct {
	provider provider
}

func (t resourceLanguageProfileType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "LanguageProfile resource",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "ID of languageprofile",
				Computed:            true,
				Type:                types.Int64Type,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				MarkdownDescription: "Name of languageprofile",
				Required:            true,
				Type:                types.StringType,
			},
			"upgrade_allowed": {
				MarkdownDescription: "Upgrade allowed Flag",
				Optional:            true,
				Computed:            true,
				Type:                types.BoolType,
			},
			//TODO: add language name validation
			"cutoff_language": {
				MarkdownDescription: "Name of language",
				Required:            true,
				Type:                types.StringType,
			},
			"languages": {
				MarkdownDescription: "list of languages in profile",
				Required:            true,
				Type:                types.SetType{ElemType: types.StringType},
			},
		},
	}, nil
}

func (t resourceLanguageProfileType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceLanguageProfile{
		provider: provider,
	}, diags
}

func (r resourceLanguageProfile) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	// Retrieve values from plan
	var plan LanguageProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build Create resource
	data := readLanguageProfile(&plan)

	// Create new LanguageProfile
	response, err := r.provider.client.AddLanguageProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create languageprofile, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "created languageprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := *writeLanguageProfile(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLanguageProfile) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	// Get current state
	var state LanguageProfile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get languageprofile current value
	response, err := r.provider.client.GetLanguageProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read languageprofiles, got error: %s", err))
		return
	}
	// Map response body to resource schema attribute
	result := *writeLanguageProfile(response)

	diags = resp.State.Set(ctx, &result)
	resp.Diagnostics.Append(diags...)
}

func (r resourceLanguageProfile) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan LanguageProfile
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get state values
	var state LanguageProfile
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build Update resource
	plan.ID.Value = state.ID.Value
	data := readLanguageProfile(&plan)

	// Update LanguageProfile
	response, err := r.provider.client.UpdateLanguageProfileContext(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update languageprofile, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "update languageprofile: "+strconv.Itoa(int(response.ID)))

	// Generate resource state struct
	result := *writeLanguageProfile(response)

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceLanguageProfile) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state LanguageProfile

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete languageprofile current value
	err := r.provider.client.DeleteLanguageProfileContext(ctx, int(state.ID.Value))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read languageprofiles, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceLanguageProfile) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	//tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ID. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("id"), id)...)
}

func writeLanguageProfile(profile *sonarr.LanguageProfile) *LanguageProfile {
	languages := make([]types.String, len(profile.Languages))
	for i, l := range profile.Languages {
		languages[i] = types.String{Value: l.Language.Name}
	}
	return &LanguageProfile{
		UpgradeAllowed: types.Bool{Value: profile.UpgradeAllowed},
		ID:             types.Int64{Value: profile.ID},
		Name:           types.String{Value: profile.Name},
		CutoffLanguage: types.String{Value: profile.Cutoff.Name},
		Languages:      languages,
	}
}

func readLanguageProfile(profile *LanguageProfile) *sonarr.LanguageProfile {
	languages := make([]sonarr.Language, len(profile.Languages))
	for i, l := range profile.Languages {
		languages[i] = sonarr.Language{
			Allowed: true,
			Language: &starr.Value{
				Name: l.Value,
				ID:   getLanguageID(l.Value),
			},
		}
	}

	return &sonarr.LanguageProfile{
		Name:           profile.Name.Value,
		UpgradeAllowed: profile.UpgradeAllowed.Value,
		Cutoff: &starr.Value{
			Name: profile.CutoffLanguage.Value,
			ID:   getLanguageID(profile.CutoffLanguage.Value),
		},
		Languages: languages,
		ID:        profile.ID.Value,
	}
}

//TODO: find a better way to maintain this find  https://github.com/Sonarr/Sonarr/blob/57e3bd8b4da16bea54a23e11d614ebd6809e2f93/src/NzbDrone.Core/Languages/Language.cs
func getLanguageID(language string) int64 {
	languages := []starr.Value{
		{
			ID:   0,
			Name: "Unknown",
		},
		{
			ID:   1,
			Name: "English",
		},
		{
			ID:   2,
			Name: "French",
		},
		{
			ID:   3,
			Name: "Spanish",
		},
		{
			ID:   4,
			Name: "German",
		},
		{
			ID:   5,
			Name: "Italian",
		},
		{
			ID:   6,
			Name: "Danish",
		},
		{
			ID:   7,
			Name: "Dutch",
		},
		{
			ID:   8,
			Name: "Japanese",
		},
		{
			ID:   9,
			Name: "Icelandic",
		},
		{
			ID:   10,
			Name: "Chinese",
		},
		{
			ID:   11,
			Name: "Russian",
		},
		{
			ID:   12,
			Name: "Polish",
		},
		{
			ID:   13,
			Name: "Vietnamese",
		},
		{
			ID:   14,
			Name: "Swedish",
		},
		{
			ID:   15,
			Name: "Norwegian",
		},
		{
			ID:   16,
			Name: "Finnish",
		},
		{
			ID:   17,
			Name: "Turkish",
		},
		{
			ID:   18,
			Name: "Portuguese",
		},
		{
			ID:   19,
			Name: "Flemish",
		},
		{
			ID:   20,
			Name: "Greek",
		},
		{
			ID:   21,
			Name: "Korean",
		},
		{
			ID:   22,
			Name: "Hungarian",
		},
		{
			ID:   23,
			Name: "Hebrew",
		},
		{
			ID:   24,
			Name: "Lithuanian",
		},
		{
			ID:   25,
			Name: "Czech",
		},
		{
			ID:   26,
			Name: "Arabic",
		},
		{
			ID:   27,
			Name: "Hindi",
		},
		{
			ID:   28,
			Name: "Bulgarian",
		},
	}
	for _, l := range languages {
		if l.Name == language {
			return l.ID
		}
	}
	return 0
}
