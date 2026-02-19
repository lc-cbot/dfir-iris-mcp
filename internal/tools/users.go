package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerUsers(s *mcp.Server, c *client.Client) {
	// List users
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_users_list",
		Description: "List all users in DFIR-IRIS",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/manage/users/list", nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get user
	type usersGetArgs struct {
		UserID int `json:"user_id" jsonschema:"User ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_users_get",
		Description: "Get details of a specific user",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args usersGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/users/%d", args.UserID)
		data, err := c.Get(ctx, path, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add user
	type usersAddArgs struct {
		UserName     string  `json:"user_name" jsonschema:"Full name of the user"`
		UserLogin    string  `json:"user_login" jsonschema:"Login username"`
		UserEmail    string  `json:"user_email" jsonschema:"Email address"`
		UserPassword string  `json:"user_password" jsonschema:"Password for the user"`
		UserIsAdmin  *bool   `json:"user_isadmin,omitempty" jsonschema:"Whether the user is an admin"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_users_add",
		Description: "Create a new user (admin operation)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args usersAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/manage/users/add", nil, toBody(args))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update user
	type usersUpdateArgs struct {
		UserID       int     `json:"user_id" jsonschema:"User ID to update"`
		UserName     *string `json:"user_name,omitempty" jsonschema:"New full name"`
		UserEmail    *string `json:"user_email,omitempty" jsonschema:"New email address"`
		UserPassword *string `json:"user_password,omitempty" jsonschema:"New password"`
		UserIsAdmin  *bool   `json:"user_isadmin,omitempty" jsonschema:"New admin status"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_users_update",
		Description: "Update a user (admin operation)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args usersUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/users/update/%d", args.UserID)
		data, err := c.Post(ctx, path, nil, toBody(args, "user_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete user
	type usersDeleteArgs struct {
		UserID int `json:"user_id" jsonschema:"User ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_users_delete",
		Description: "Delete a user (admin operation, irreversible)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args usersDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/users/delete/%d", args.UserID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
