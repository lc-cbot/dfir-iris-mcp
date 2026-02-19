package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerEvidences(s *mcp.Server, c *client.Client) {
	// List evidences
	type evidencesListArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_evidences_list",
		Description: "List all evidences in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args evidencesListArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/case/evidences/list", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get evidence
	type evidencesGetArgs struct {
		CaseID     int `json:"case_id" jsonschema:"Case ID"`
		EvidenceID int `json:"evidence_id" jsonschema:"Evidence ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_evidences_get",
		Description: "Get details of a specific evidence item",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args evidencesGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/evidences/%d", args.EvidenceID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add evidence
	type evidencesAddArgs struct {
		CaseID          int     `json:"case_id" jsonschema:"Case ID"`
		Filename        string  `json:"filename" jsonschema:"Filename of the evidence"`
		FileSize        *int    `json:"file_size,omitempty" jsonschema:"File size in bytes"`
		FileHash        *string `json:"file_hash,omitempty" jsonschema:"File hash (MD5, SHA1, or SHA256)"`
		FileDescription *string `json:"file_description,omitempty" jsonschema:"Description of the evidence"`
		EvidenceTypeID  *int    `json:"type_id,omitempty" jsonschema:"Evidence type ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_evidences_add",
		Description: "Add a new evidence record to a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args evidencesAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/case/evidences/add", cidQuery(args.CaseID), toBody(args, "case_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update evidence
	type evidencesUpdateArgs struct {
		CaseID          int     `json:"case_id" jsonschema:"Case ID"`
		EvidenceID      int     `json:"evidence_id" jsonschema:"Evidence ID to update"`
		Filename        *string `json:"filename,omitempty" jsonschema:"New filename"`
		FileSize        *int    `json:"file_size,omitempty" jsonschema:"New file size"`
		FileHash        *string `json:"file_hash,omitempty" jsonschema:"New file hash"`
		FileDescription *string `json:"file_description,omitempty" jsonschema:"New description"`
		EvidenceTypeID  *int    `json:"type_id,omitempty" jsonschema:"New evidence type ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_evidences_update",
		Description: "Update an evidence record in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args evidencesUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/evidences/update/%d", args.EvidenceID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "evidence_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete evidence
	type evidencesDeleteArgs struct {
		CaseID     int `json:"case_id" jsonschema:"Case ID"`
		EvidenceID int `json:"evidence_id" jsonschema:"Evidence ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_evidences_delete",
		Description: "Delete an evidence record from a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args evidencesDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/evidences/delete/%d", args.EvidenceID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
