package tools

import (
	"context"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerSystem(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_system_ping",
		Description: "Ping the DFIR-IRIS server to check connectivity",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/api/ping", nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_system_versions",
		Description: "Get DFIR-IRIS server version information",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/api/versions", nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
