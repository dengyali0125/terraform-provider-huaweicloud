package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscExportTemplate_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_export_template.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDSCScanTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscExportTemplate_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "bucket"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "end_time"),
				),
			},
		},
	})
}

func testAccDataSourceDscExportTemplate_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_export_template" "test" {
  template_id = "%s"
  bucket      = "dyl-test1"
}
`, acceptance.HW_DSC_SCAN_TEMPLATE_ID)
}
