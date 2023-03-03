package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAutoTagConditionRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionRootFolderDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sonarr_auto_tag_condition_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.sonarr_auto_tag_condition_root_folder.test", "name", "Test"),
					resource.TestCheckResourceAttr("sonarr_auto_tag.test", "specifications.0.value", "/config")),
			},
		},
	})
}

const testAccAutoTagConditionRootFolderDataSourceConfig = `
resource "sonarr_tag" "test" {
	label = "atconditionfolder"
}

data  "sonarr_auto_tag_condition_root_folder" "test" {
	name = "Test"
	negate = false
	required = false
	value = "/config"
}

resource "sonarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDSRootFolder"

	tags = [sonarr_tag.test.id]
	
	specifications = [data.sonarr_auto_tag_condition_root_folder.test]	
}`