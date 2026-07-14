package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC POST /v1/{project_id}/scan-templates/{template_id}/export
func DataSourceDscExportTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscExportTemplateRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The template ID.",
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The OBS bucket name.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The export status.",
			},
			"full_file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full path file name.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time.",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The start time.",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The end time.",
			},
			"failed_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The failed description.",
			},
		},
	}
}

func buildDscExportTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{}

	if v, ok := d.GetOk("bucket"); ok {
		bodyParams["bucket"] = v
	}

	return bodyParams
}

func dataSourceDscExportTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/scan-templates/{template_id}/export"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", d.Get("template_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
		JSONBody: buildDscExportTemplateBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error exporting DSC template: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("bucket", utils.PathSearch("bucket", respBody, nil)),
		d.Set("full_file_name", utils.PathSearch("full_file_name", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("start_time", utils.PathSearch("start_time", respBody, nil)),
		d.Set("end_time", utils.PathSearch("end_time", respBody, nil)),
		d.Set("failed_description", utils.PathSearch("failed_description", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
