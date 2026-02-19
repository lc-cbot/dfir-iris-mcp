package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerGroups(s *mcp.Server, c *client.Client) {
	// List groups
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_groups_list",
		Description: "List all groups in DFIR-IRIS",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/manage/groups/list", nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add group
	type groupsAddArgs struct {
		GroupName        string  `json:"group_name" jsonschema:"Name of the group"`
		GroupDescription *string `json:"group_description,omitempty" jsonschema:"Group description"`
		GroupPermissions *int    `json:"group_permissions,omitempty" jsonschema:"Permission bitmask"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_groups_add",
		Description: "Create a new group (admin operation)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args groupsAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/manage/groups/add", nil, toBody(args))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update group
	type groupsUpdateArgs struct {
		GroupID          int     `json:"group_id" jsonschema:"Group ID to update"`
		GroupName        *string `json:"group_name,omitempty" jsonschema:"New group name"`
		GroupDescription *string `json:"group_description,omitempty" jsonschema:"New description"`
		GroupPermissions *int    `json:"group_permissions,omitempty" jsonschema:"New permission bitmask"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_groups_update",
		Description: "Update a group (admin operation)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args groupsUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/groups/update/%d", args.GroupID)
		data, err := c.Post(ctx, path, nil, toBody(args, "group_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete group
	type groupsDeleteArgs struct {
		GroupID int `json:"group_id" jsonschema:"Group ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_groups_delete",
		Description: "Delete a group (admin operation, irreversible)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args groupsDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/groups/delete/%d", args.GroupID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
