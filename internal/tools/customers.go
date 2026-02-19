package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerCustomers(s *mcp.Server, c *client.Client) {
	// List customers
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_customers_list",
		Description: "List all customers in DFIR-IRIS",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/manage/customers/list", nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add customer
	type customersAddArgs struct {
		CustomerName        string  `json:"customer_name" jsonschema:"Customer name"`
		CustomerDescription *string `json:"customer_description,omitempty" jsonschema:"Customer description"`
		CustomerSLA         *string `json:"customer_sla,omitempty" jsonschema:"SLA terms"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_customers_add",
		Description: "Create a new customer",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args customersAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/manage/customers/add", nil, toBody(args))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update customer
	type customersUpdateArgs struct {
		CustomerID          int     `json:"customer_id" jsonschema:"Customer ID to update"`
		CustomerName        *string `json:"customer_name,omitempty" jsonschema:"New customer name"`
		CustomerDescription *string `json:"customer_description,omitempty" jsonschema:"New description"`
		CustomerSLA         *string `json:"customer_sla,omitempty" jsonschema:"New SLA terms"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_customers_update",
		Description: "Update a customer",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args customersUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/customers/update/%d", args.CustomerID)
		data, err := c.Post(ctx, path, nil, toBody(args, "customer_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete customer
	type customersDeleteArgs struct {
		CustomerID int `json:"customer_id" jsonschema:"Customer ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_customers_delete",
		Description: "Delete a customer (irreversible)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args customersDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/manage/customers/delete/%d", args.CustomerID)
		data, err := c.Post(ctx, path, nil, nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
