package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerAssets(s *mcp.Server, c *client.Client) {
	// List assets
	type assetsListArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_assets_list",
		Description: "List all assets in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args assetsListArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/case/assets/list", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get asset
	type assetsGetArgs struct {
		CaseID  int `json:"case_id" jsonschema:"Case ID"`
		AssetID int `json:"asset_id" jsonschema:"Asset ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_assets_get",
		Description: "Get details of a specific asset in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args assetsGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/assets/%d", args.AssetID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add asset
	type assetsAddArgs struct {
		CaseID           int     `json:"case_id" jsonschema:"Case ID"`
		AssetName        string  `json:"asset_name" jsonschema:"Name of the asset (e.g. hostname or IP)"`
		AssetTypeID      int     `json:"asset_type_id" jsonschema:"Asset type ID (use settings_asset_types to list)"`
		AssetDescription *string `json:"asset_description,omitempty" jsonschema:"Description of the asset"`
		AssetIP          *string `json:"asset_ip,omitempty" jsonschema:"IP address of the asset"`
		AssetDomain      *string `json:"asset_domain,omitempty" jsonschema:"Domain of the asset"`
		AssetTags        *string                 `json:"asset_tags,omitempty" jsonschema:"Comma-separated tags"`
		AnalysisStatus   *int                    `json:"analysis_status,omitempty" jsonschema:"Analysis status ID"`
		CompromiseStatus *int                    `json:"compromise_status_id,omitempty" jsonschema:"Compromise status ID"`
		CustomAttributes *map[string]interface{} `json:"custom_attributes,omitempty" jsonschema:"Custom attributes as key-value pairs (e.g. {\"limacharlie_sid\": \"uuid\"})"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_assets_add",
		Description: "Add a new asset to a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args assetsAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/case/assets/add", cidQuery(args.CaseID), toBody(args, "case_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update asset
	type assetsUpdateArgs struct {
		CaseID           int     `json:"case_id" jsonschema:"Case ID"`
		AssetID          int     `json:"asset_id" jsonschema:"Asset ID to update"`
		AssetName        *string `json:"asset_name,omitempty" jsonschema:"New asset name"`
		AssetTypeID      *int    `json:"asset_type_id,omitempty" jsonschema:"New asset type ID"`
		AssetDescription *string `json:"asset_description,omitempty" jsonschema:"New description"`
		AssetIP          *string `json:"asset_ip,omitempty" jsonschema:"New IP address"`
		AssetDomain      *string `json:"asset_domain,omitempty" jsonschema:"New domain"`
		AssetTags        *string                 `json:"asset_tags,omitempty" jsonschema:"New comma-separated tags"`
		AnalysisStatus   *int                    `json:"analysis_status,omitempty" jsonschema:"New analysis status ID"`
		CompromiseStatus *int                    `json:"compromise_status_id,omitempty" jsonschema:"New compromise status ID"`
		CustomAttributes *map[string]interface{} `json:"custom_attributes,omitempty" jsonschema:"Custom attributes as key-value pairs"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_assets_update",
		Description: "Update an existing asset in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args assetsUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/assets/update/%d", args.AssetID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "asset_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete asset
	type assetsDeleteArgs struct {
		CaseID  int `json:"case_id" jsonschema:"Case ID"`
		AssetID int `json:"asset_id" jsonschema:"Asset ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_assets_delete",
		Description: "Delete an asset from a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args assetsDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/assets/delete/%d", args.AssetID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
