---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_export_template"
description: |-
  Use this data source to export the template content within HuaweiCloud.
---

# huaweicloud_dsc_export_template

Use this data source to export the template content within HuaweiCloud.

## Example Usage

```hcl
variable "template_id" {}

data "huaweicloud_dsc_export_template" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `template_id` - (Required, String) Specifies the template ID.

* `bucket` - (Required, String) Specifies the OBS bucket name for exporting the template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `status` - The export status.

* `full_file_name` - The full path file name.

* `create_time` - The creation time.

* `start_time` - The start time.

* `end_time` - The end time.

* `failed_description` - The failed description.
