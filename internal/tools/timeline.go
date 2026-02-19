package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTimeline(s *mcp.Server, c *client.Client) {
	// List timeline events
	type timelineListArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_timeline_list",
		Description: "List all timeline events in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args timelineListArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/case/timeline/events/list", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get timeline event
	type timelineGetArgs struct {
		CaseID  int `json:"case_id" jsonschema:"Case ID"`
		EventID int `json:"event_id" jsonschema:"Event ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_timeline_get",
		Description: "Get details of a specific timeline event",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args timelineGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/timeline/events/%d", args.EventID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add timeline event
	type timelineAddArgs struct {
		CaseID          int     `json:"case_id" jsonschema:"Case ID"`
		EventTitle      string  `json:"event_title" jsonschema:"Title of the event"`
		EventDate       string  `json:"event_date" jsonschema:"Date/time of the event (format: YYYY-MM-DDTHH:MM:SS.000)"`
		EventTZ         string  `json:"event_tz" jsonschema:"Timezone offset (e.g. +00:00, -05:00, +02:00)"`
		EventCategoryID int     `json:"event_category_id" jsonschema:"Event category ID (use settings_event_categories to list)"`
		EventAssets     []int   `json:"event_assets" jsonschema:"List of asset IDs linked to this event (use empty list [] if none)"`
		EventIOCs       []int   `json:"event_iocs" jsonschema:"List of IOC IDs linked to this event (use empty list [] if none)"`
		EventContent    *string `json:"event_content,omitempty" jsonschema:"Event content/description"`
		EventRaw        *string `json:"event_raw,omitempty" jsonschema:"Raw event data"`
		EventSource     *string `json:"event_source,omitempty" jsonschema:"Source of the event"`
		EventColor      *string `json:"event_color,omitempty" jsonschema:"Color hex code for display"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_timeline_add",
		Description: "Add a new event to the case timeline",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args timelineAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/case/timeline/events/add", cidQuery(args.CaseID), toBody(args, "case_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update timeline event
	type timelineUpdateArgs struct {
		CaseID          int     `json:"case_id" jsonschema:"Case ID"`
		EventID         int     `json:"event_id" jsonschema:"Event ID to update"`
		EventTitle      *string `json:"event_title,omitempty" jsonschema:"New event title"`
		EventDate       *string `json:"event_date,omitempty" jsonschema:"New date/time"`
		EventContent    *string `json:"event_content,omitempty" jsonschema:"New content"`
		EventRaw        *string `json:"event_raw,omitempty" jsonschema:"New raw data"`
		EventSource     *string `json:"event_source,omitempty" jsonschema:"New source"`
		EventCategoryID *int    `json:"event_category_id,omitempty" jsonschema:"New category ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_timeline_update",
		Description: "Update a timeline event in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args timelineUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/timeline/events/update/%d", args.EventID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "event_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete timeline event
	type timelineDeleteArgs struct {
		CaseID  int `json:"case_id" jsonschema:"Case ID"`
		EventID int `json:"event_id" jsonschema:"Event ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_timeline_delete",
		Description: "Delete a timeline event from a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args timelineDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/timeline/events/delete/%d", args.EventID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
