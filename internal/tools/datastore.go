package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerDatastore(s *mcp.Server, c *client.Client) {
	// List datastore tree
	type datastoreTreeArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_tree",
		Description: "Get the datastore folder/file tree for a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreTreeArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/datastore/list/tree", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get file info
	type datastoreFileGetArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		FileID int `json:"file_id" jsonschema:"Datastore file ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_file_get",
		Description: "Get metadata of a file in the datastore",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFileGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/file/info/%d", args.FileID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add file (metadata only — binary upload not supported via MCP)
	type datastoreFileAddArgs struct {
		CaseID           int     `json:"case_id" jsonschema:"Case ID"`
		ParentID         int     `json:"parent_id" jsonschema:"Parent folder ID"`
		FileOriginalName string  `json:"file_original_name" jsonschema:"Original filename"`
		FileDescription  *string `json:"file_description,omitempty" jsonschema:"File description"`
		FilePassword     *string `json:"file_password,omitempty" jsonschema:"Password if file is encrypted"`
		FileIsIoc        *bool   `json:"file_is_ioc,omitempty" jsonschema:"Whether file is an IOC"`
		FileIsEvidence   *bool   `json:"file_is_evidence,omitempty" jsonschema:"Whether file is evidence"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_file_add",
		Description: "Add a file entry to the datastore (metadata only, binary upload not supported via MCP)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFileAddArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/file/add/%d", args.ParentID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "parent_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update file
	type datastoreFileUpdateArgs struct {
		CaseID           int     `json:"case_id" jsonschema:"Case ID"`
		FileID           int     `json:"file_id" jsonschema:"File ID to update"`
		FileOriginalName *string `json:"file_original_name,omitempty" jsonschema:"New filename"`
		FileDescription  *string `json:"file_description,omitempty" jsonschema:"New description"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_file_update",
		Description: "Update a file's metadata in the datastore",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFileUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/file/update/%d", args.FileID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "file_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete file
	type datastoreFileDeleteArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		FileID int `json:"file_id" jsonschema:"File ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_file_delete",
		Description: "Delete a file from the datastore",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFileDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/file/delete/%d", args.FileID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Move file
	type datastoreFileMoveArgs struct {
		CaseID              int `json:"case_id" jsonschema:"Case ID"`
		FileID              int `json:"file_id" jsonschema:"File ID to move"`
		DestinationFolderID int `json:"destination_folder_id" jsonschema:"Destination folder ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_file_move",
		Description: "Move a file to a different folder in the datastore",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFileMoveArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/file/move/%d", args.FileID)
		body := map[string]interface{}{"destination_folder_id": args.DestinationFolderID}
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add folder
	type datastoreFolderAddArgs struct {
		CaseID     int    `json:"case_id" jsonschema:"Case ID"`
		FolderName string `json:"folder_name" jsonschema:"Name of the new folder"`
		ParentID   int    `json:"parent_id" jsonschema:"Parent folder ID (0 for root)"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_folder_add",
		Description: "Create a new folder in the datastore",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFolderAddArgs) (*mcp.CallToolResult, any, error) {
		body := map[string]interface{}{
			"folder_name": args.FolderName,
			"parent_id":   args.ParentID,
		}
		data, err := c.Post(ctx, "/datastore/folder/add", cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete folder
	type datastoreFolderDeleteArgs struct {
		CaseID   int `json:"case_id" jsonschema:"Case ID"`
		FolderID int `json:"folder_id" jsonschema:"Folder ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_folder_delete",
		Description: "Delete a folder from the datastore (and all contents)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFolderDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/folder/delete/%d", args.FolderID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Rename folder
	type datastoreFolderRenameArgs struct {
		CaseID     int    `json:"case_id" jsonschema:"Case ID"`
		FolderID   int    `json:"folder_id" jsonschema:"Folder ID to rename"`
		FolderName string `json:"folder_name" jsonschema:"New folder name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_folder_rename",
		Description: "Rename a folder in the datastore",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFolderRenameArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/folder/rename/%d", args.FolderID)
		body := map[string]interface{}{"folder_name": args.FolderName}
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Move folder
	type datastoreFolderMoveArgs struct {
		CaseID              int `json:"case_id" jsonschema:"Case ID"`
		FolderID            int `json:"folder_id" jsonschema:"Folder ID to move"`
		DestinationFolderID int `json:"destination_folder_id" jsonschema:"Destination parent folder ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_datastore_folder_move",
		Description: "Move a folder to a different parent folder in the datastore",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args datastoreFolderMoveArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/datastore/folder/move/%d", args.FolderID)
		body := map[string]interface{}{"destination_folder_id": args.DestinationFolderID}
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
