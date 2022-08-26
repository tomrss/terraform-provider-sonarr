package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexersDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a delay profile to have a value to check
			{
				Config: testAccIndexerResourceConfig("datasourceTest", "true"),
			},
			// Read testing
			{
				Config: testAccIndexersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.sonarr_indexers.test", "indexers.*", map[string]string{"protocol": "usenet"}),
				),
			},
		},
	})
}

const testAccIndexersDataSourceConfig = `
data "sonarr_indexers" "test" {
}
`
