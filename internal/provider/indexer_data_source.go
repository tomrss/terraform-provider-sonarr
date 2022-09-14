package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golift.io/starr/sonarr"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerDataSource{}

func NewIndexerDataSource() datasource.DataSource {
	return &IndexerDataSource{}
}

// IndexerDataSource defines the indexer implementation.
type IndexerDataSource struct {
	client *sonarr.Sonarr
}

func (d *IndexerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_indexer"
}

func (d *IndexerDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "Single [Indexer](../resources/indexer).",
		Attributes: map[string]tfsdk.Attribute{
			"enable_automatic_search": {
				MarkdownDescription: "Enable automatic search flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_interactive_search": {
				MarkdownDescription: "Enable interactive search flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"enable_rss": {
				MarkdownDescription: "Enable RSS flag.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"priority": {
				MarkdownDescription: "Priority.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"download_client_id": {
				MarkdownDescription: "Download client ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"config_contract": {
				MarkdownDescription: "Indexer configuration template.",
				Computed:            true,
				Type:                types.StringType,
			},
			"implementation": {
				MarkdownDescription: "Indexer implementation name.",
				Computed:            true,
				Type:                types.StringType,
			},
			"name": {
				MarkdownDescription: "Indexer name.",
				Required:            true,
				Type:                types.StringType,
			},
			"protocol": {
				MarkdownDescription: "Protocol. Valid values are 'usenet' and 'torrent'.",
				Computed:            true,
				Type:                types.StringType,
			},
			"tags": {
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"id": {
				MarkdownDescription: "Indexer ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			// Field values
			"allow_zero_size": {
				MarkdownDescription: "Allow zero size files.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"anime_standard_format_search": {
				MarkdownDescription: "Search anime in standard format.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"ranked_only": {
				MarkdownDescription: "Allow ranked only.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"delay": {
				MarkdownDescription: "Delay before grabbing.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"minimum_seeders": {
				MarkdownDescription: "Minimum seeders.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"season_pack_seed_time": {
				MarkdownDescription: "Season seed time.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_time": {
				MarkdownDescription: "Seed time.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"seed_ratio": {
				MarkdownDescription: "Seed ratio.",
				Computed:            true,
				Type:                types.Float64Type,
			},
			"additional_parameters": {
				MarkdownDescription: "Additional parameters.",
				Computed:            true,
				Type:                types.StringType,
			},
			"api_key": {
				MarkdownDescription: "API key.",
				Computed:            true,
				Type:                types.StringType,
			},
			"api_path": {
				MarkdownDescription: "API path.",
				Computed:            true,
				Type:                types.StringType,
			},
			"base_url": {
				MarkdownDescription: "Base URL.",
				Computed:            true,
				Type:                types.StringType,
			},
			"captcha_token": {
				MarkdownDescription: "Captcha token.",
				Computed:            true,
				Type:                types.StringType,
			},
			"cookie": {
				MarkdownDescription: "Cookie.",
				Computed:            true,
				Type:                types.StringType,
			},
			"passkey": {
				MarkdownDescription: "Passkey.",
				Computed:            true,
				Type:                types.StringType,
			},
			"username": {
				MarkdownDescription: "Username.",
				Computed:            true,
				Type:                types.StringType,
			},
			"categories": {
				MarkdownDescription: "Series list.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
			"anime_categories": {
				MarkdownDescription: "Anime list.",
				Computed:            true,
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
			},
		},
	}, nil
}

func (d *IndexerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sonarr.Sonarr)
	if !ok {
		resp.Diagnostics.AddError(
			UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *sonarr.Sonarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *IndexerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Indexer

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get indexer current value
	response, err := d.client.GetIndexersContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(ClientError, fmt.Sprintf("Unable to read indexer, got error: %s", err))

		return
	}

	indexer, err := findIndexer(data.Name.Value, response)
	if err != nil {
		resp.Diagnostics.AddError(DataSourceError, fmt.Sprintf("Unable to find indexer, got error: %s", err))

		return
	}

	result := writeIndexer(ctx, indexer)
	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func findIndexer(name string, indexers []*sonarr.IndexerOutput) (*sonarr.IndexerOutput, error) {
	for _, i := range indexers {
		if i.Name == name {
			return i, nil
		}
	}

	return nil, fmt.Errorf("no language indexer with name %s", name)
}
