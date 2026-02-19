package tools

import (
	"context"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerSettings(s *mcp.Server, c *client.Client) {
	settings := []struct {
		name string
		desc string
		path string
	}{
		{"dfir_iris_settings_asset_types", "List available asset types for use when creating assets", "/manage/asset-type/list"},
		{"dfir_iris_settings_ioc_types", "List available IOC types for use when creating IOCs", "/manage/ioc-types/list"},
		{"dfir_iris_settings_task_status", "List available task statuses", "/manage/task-status/list"},
		{"dfir_iris_settings_analysis_status", "List available analysis statuses", "/manage/analysis-status/list"},
		{"dfir_iris_settings_case_states", "List available case states", "/manage/case-states/list"},
		{"dfir_iris_settings_case_templates", "List available case templates", "/manage/case-templates/list"},
		{"dfir_iris_settings_classifications", "List available case classifications", "/manage/case-classifications/list"},
		{"dfir_iris_settings_evidence_types", "List available evidence types", "/manage/evidence-types/list"},
		{"dfir_iris_settings_event_categories", "List available event categories for timeline events", "/manage/event-categories/list"},
	}

	for _, st := range settings {
		st := st
		mcp.AddTool(s, &mcp.Tool{
			Name:        st.name,
			Description: st.desc,
		}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
			data, err := c.Get(ctx, st.path, nil)
			if err != nil {
				return errorResult(err), nil, nil
			}
			return textResult(data), nil, nil
		})
	}
}
