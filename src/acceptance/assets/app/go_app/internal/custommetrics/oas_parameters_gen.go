// Code generated by ogen, DO NOT EDIT.

package api

// V1AppsAppGuidMetricsPostParams is parameters of POST /v1/apps/{appGuid}/metrics operation.
type V1AppsAppGuidMetricsPostParams struct {
	// The GUID identifying the application the custom metrics are submitted for. Can be found in the
	// `application_id` property of the JSON object stored in the `VCAP_APPLICATION` environment variable.
	AppGuid string
}
