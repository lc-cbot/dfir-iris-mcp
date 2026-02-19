package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerCases(s *mcp.Server, c *client.Client) {
	// List all cases
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_list",
		Description: "List all cases in DFIR-IRIS",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/manage/cases/list", nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Filter cases
	type casesFilterArgs struct {
		CaseName     *string `json:"case_name,omitempty" jsonschema:"Filter by case name (substring match)"`
		CaseCustomer *int    `json:"case_customer,omitempty" jsonschema:"Filter by customer ID"`
		CaseState    *int    `json:"case_state,omitempty" jsonschema:"Filter by case state ID"`
		Page         *int    `json:"page,omitempty" jsonschema:"Page number for pagination"`
		PerPage      *int    `json:"per_page,omitempty" jsonschema:"Results per page"`
		Sort         *string `json:"sort,omitempty" jsonschema:"Sort field"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_filter",
		Description: "Filter cases with optional search criteria",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesFilterArgs) (*mcp.CallToolResult, any, error) {
		q := toQuery(args)
		data, err := c.Get(ctx, "/manage/cases/filter", q)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add a case
	type casesAddArgs struct {
		CaseName       string  `json:"case_name" jsonschema:"Name of the case"`
		CaseCustomer   int     `json:"case_customer" jsonschema:"Customer ID to associate with the case"`
		CaseDescription *string `json:"case_description,omitempty" jsonschema:"Case description"`
		CaseSOCID      *string `json:"case_soc_id,omitempty" jsonschema:"SOC ticket ID"`
		ClassificationID *int  `json:"classification_id,omitempty" jsonschema:"Classification ID"`
		CaseTemplateID *int    `json:"case_template_id,omitempty" jsonschema:"Case template ID to apply"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_add",
		Description: "Create a new case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/manage/cases/add", nil, toBody(args))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update a case
	type casesUpdateArgs struct {
		CaseID          int     `json:"case_id" jsonschema:"ID of the case to update"`
		CaseName        *string `json:"case_name,omitempty" jsonschema:"New case name"`
		CaseDescription *string `json:"case_description,omitempty" jsonschema:"New case description"`
		CaseCustomer    *int    `json:"case_customer,omitempty" jsonschema:"New customer ID"`
		CaseSOCID       *string `json:"case_soc_id,omitempty" jsonschema:"New SOC ticket ID"`
		ClassificationID *int   `json:"classification_id,omitempty" jsonschema:"New classification ID"`
		StateID         *int    `json:"state_id,omitempty" jsonschema:"New state ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_update",
		Description: "Update an existing case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/cases/update/%d", args.CaseID)
		data, err := c.Post(ctx, path, nil, toBody(args, "case_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete a case
	type casesDeleteArgs struct {
		CaseID int `json:"case_id" jsonschema:"ID of the case to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_delete",
		Description: "Delete a case (irreversible)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/cases/delete/%d", args.CaseID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Close a case
	type casesCloseArgs struct {
		CaseID int `json:"case_id" jsonschema:"ID of the case to close"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_close",
		Description: "Close a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesCloseArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/cases/close/%d", args.CaseID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Reopen a case
	type casesReopenArgs struct {
		CaseID int `json:"case_id" jsonschema:"ID of the case to reopen"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_reopen",
		Description: "Reopen a previously closed case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesReopenArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/cases/reopen/%d", args.CaseID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update case summary
	type casesSummaryUpdateArgs struct {
		CaseID      int    `json:"case_id" jsonschema:"Case ID"`
		CaseSummary string `json:"case_summary" jsonschema:"New case summary text (supports markdown)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_summary_update",
		Description: "Update the summary/description of a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesSummaryUpdateArgs) (*mcp.CallToolResult, any, error) {
		body := map[string]interface{}{"case_summary": args.CaseSummary}
		data, err := c.Post(ctx, "/case/summary/update", cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Export a case
	type casesExportArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID to export"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_cases_export",
		Description: "Export a case as JSON",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args casesExportArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/case/export", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
