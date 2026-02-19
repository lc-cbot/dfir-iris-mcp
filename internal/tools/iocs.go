package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerIOCs(s *mcp.Server, c *client.Client) {
	// List IOCs
	type iocsListArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_iocs_list",
		Description: "List all IOCs in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args iocsListArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/case/ioc/list", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get IOC
	type iocsGetArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		IOCID  int `json:"ioc_id" jsonschema:"IOC ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_iocs_get",
		Description: "Get details of a specific IOC in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args iocsGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/ioc/%d", args.IOCID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add IOC
	type iocsAddArgs struct {
		CaseID         int     `json:"case_id" jsonschema:"Case ID"`
		IOCValue       string  `json:"ioc_value" jsonschema:"IOC value (e.g. IP, hash, domain)"`
		IOCTypeID      int     `json:"ioc_type_id" jsonschema:"IOC type ID (use settings_ioc_types to list)"`
		IOCDescription *string `json:"ioc_description,omitempty" jsonschema:"Description of the IOC"`
		IOCTLPID       *int    `json:"ioc_tlp_id,omitempty" jsonschema:"TLP level ID"`
		IOCTags        *string `json:"ioc_tags,omitempty" jsonschema:"Comma-separated tags"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_iocs_add",
		Description: "Add a new IOC to a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args iocsAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/case/ioc/add", cidQuery(args.CaseID), toBody(args, "case_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update IOC
	type iocsUpdateArgs struct {
		CaseID         int     `json:"case_id" jsonschema:"Case ID"`
		IOCID          int     `json:"ioc_id" jsonschema:"IOC ID to update"`
		IOCValue       *string `json:"ioc_value,omitempty" jsonschema:"New IOC value"`
		IOCTypeID      *int    `json:"ioc_type_id,omitempty" jsonschema:"New IOC type ID"`
		IOCDescription *string `json:"ioc_description,omitempty" jsonschema:"New description"`
		IOCTLPID       *int    `json:"ioc_tlp_id,omitempty" jsonschema:"New TLP level ID"`
		IOCTags        *string `json:"ioc_tags,omitempty" jsonschema:"New comma-separated tags"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_iocs_update",
		Description: "Update an existing IOC in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args iocsUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/ioc/update/%d", args.IOCID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "ioc_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete IOC
	type iocsDeleteArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		IOCID  int `json:"ioc_id" jsonschema:"IOC ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_iocs_delete",
		Description: "Delete an IOC from a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args iocsDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/ioc/delete/%d", args.IOCID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
