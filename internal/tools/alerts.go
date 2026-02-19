package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerAlerts(s *mcp.Server, c *client.Client) {
	// Filter alerts
	type alertsFilterArgs struct {
		AlertTitle          *string `json:"alert_title,omitempty" jsonschema:"Filter by alert title"`
		AlertSeverityID     *int    `json:"alert_severity_id,omitempty" jsonschema:"Filter by severity ID"`
		AlertStatusID       *int    `json:"alert_status_id,omitempty" jsonschema:"Filter by status ID"`
		AlertCustomerID     *int    `json:"alert_customer_id,omitempty" jsonschema:"Filter by customer ID"`
		AlertSource         *string `json:"alert_source,omitempty" jsonschema:"Filter by alert source"`
		AlertClassificationID *int  `json:"alert_classification_id,omitempty" jsonschema:"Filter by classification ID"`
		Page                *int    `json:"page,omitempty" jsonschema:"Page number"`
		PerPage             *int    `json:"per_page,omitempty" jsonschema:"Results per page"`
		Sort                *string `json:"sort,omitempty" jsonschema:"Sort field"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_filter",
		Description: "Filter alerts with optional search criteria",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsFilterArgs) (*mcp.CallToolResult, any, error) {
		q := toQuery(args)
		data, err := c.Get(ctx, "/alerts/filter", q)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get a single alert
	type alertsGetArgs struct {
		AlertID int `json:"alert_id" jsonschema:"Alert ID to retrieve"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_get",
		Description: "Get details of a specific alert",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/alerts/%d", args.AlertID)
		data, err := c.Get(ctx, path, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add an alert
	type alertsAddArgs struct {
		AlertTitle          string  `json:"alert_title" jsonschema:"Title of the alert"`
		AlertSeverityID     int     `json:"alert_severity_id" jsonschema:"Severity ID"`
		AlertStatusID       int     `json:"alert_status_id" jsonschema:"Status ID"`
		AlertCustomerID     int     `json:"alert_customer_id" jsonschema:"Customer ID"`
		AlertDescription    *string `json:"alert_description,omitempty" jsonschema:"Alert description"`
		AlertSource         *string `json:"alert_source,omitempty" jsonschema:"Source of the alert (e.g. SIEM name)"`
		AlertSourceRef      *string `json:"alert_source_ref,omitempty" jsonschema:"Source reference ID"`
		AlertSourceLink     *string `json:"alert_source_link,omitempty" jsonschema:"Link to alert in source system"`
		AlertClassificationID *int  `json:"alert_classification_id,omitempty" jsonschema:"Classification ID"`
		AlertNote           *string `json:"alert_note,omitempty" jsonschema:"Alert note"`
		AlertTags           *string `json:"alert_tags,omitempty" jsonschema:"Comma-separated tags"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_add",
		Description: "Create a new alert",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/alerts/add", nil, toBody(args))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update an alert
	type alertsUpdateArgs struct {
		AlertID             int     `json:"alert_id" jsonschema:"ID of the alert to update"`
		AlertTitle          *string `json:"alert_title,omitempty" jsonschema:"New alert title"`
		AlertSeverityID     *int    `json:"alert_severity_id,omitempty" jsonschema:"New severity ID"`
		AlertStatusID       *int    `json:"alert_status_id,omitempty" jsonschema:"New status ID"`
		AlertCustomerID     *int    `json:"alert_customer_id,omitempty" jsonschema:"New customer ID"`
		AlertDescription    *string `json:"alert_description,omitempty" jsonschema:"New description"`
		AlertSource         *string `json:"alert_source,omitempty" jsonschema:"New source"`
		AlertClassificationID *int  `json:"alert_classification_id,omitempty" jsonschema:"New classification ID"`
		AlertNote           *string `json:"alert_note,omitempty" jsonschema:"New note"`
		AlertTags           *string `json:"alert_tags,omitempty" jsonschema:"New comma-separated tags"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_update",
		Description: "Update an existing alert",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/alerts/update/%d", args.AlertID)
		data, err := c.Post(ctx, path, nil, toBody(args, "alert_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete an alert
	type alertsDeleteArgs struct {
		AlertID int `json:"alert_id" jsonschema:"ID of the alert to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_delete",
		Description: "Delete an alert (irreversible)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/alerts/delete/%d", args.AlertID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Escalate an alert to a case
	type alertsEscalateArgs struct {
		AlertID        int   `json:"alert_id" jsonschema:"ID of the alert to escalate"`
		IOCsImport     *bool `json:"iocs_import,omitempty" jsonschema:"Import IOCs from the alert into the case"`
		AssetsImport   *bool `json:"assets_import,omitempty" jsonschema:"Import assets from the alert into the case"`
		CaseID         *int  `json:"case_id,omitempty" jsonschema:"Existing case ID to escalate into (creates new case if omitted)"`
		CaseTemplateID *int  `json:"case_template_id,omitempty" jsonschema:"Case template ID for the new case"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_escalate",
		Description: "Escalate an alert to a new or existing case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsEscalateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/alerts/escalate/%d", args.AlertID)
		data, err := c.Post(ctx, path, nil, toBody(args, "alert_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Merge an alert into a case
	type alertsMergeArgs struct {
		AlertID      int `json:"alert_id" jsonschema:"ID of the alert to merge"`
		TargetCaseID int `json:"target_case_id" jsonschema:"Case ID to merge the alert into"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_merge",
		Description: "Merge an alert into an existing case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsMergeArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/alerts/merge/%d", args.AlertID)
		body := map[string]interface{}{"target_case_id": args.TargetCaseID}
		data, err := c.Post(ctx, path, nil, body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Unmerge an alert from a case
	type alertsUnmergeArgs struct {
		AlertID int `json:"alert_id" jsonschema:"ID of the alert to unmerge from its case"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_alerts_unmerge",
		Description: "Unmerge an alert from its associated case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args alertsUnmergeArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/alerts/unmerge/%d", args.AlertID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
