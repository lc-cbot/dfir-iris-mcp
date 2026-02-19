package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerComments(s *mcp.Server, c *client.Client) {
	// List comments
	type commentsListArgs struct {
		CaseID     int    `json:"case_id" jsonschema:"Case ID"`
		ObjectType string `json:"object_type" jsonschema:"Object type (e.g. cases, assets, ioc, timeline_events, tasks, evidences)"`
		ObjectID   int    `json:"object_id" jsonschema:"Object ID to list comments for"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_comments_list",
		Description: "List comments on a case object (asset, IOC, event, task, etc.)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args commentsListArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/%s/%d/comments/list", args.ObjectType, args.ObjectID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add comment
	type commentsAddArgs struct {
		CaseID      int    `json:"case_id" jsonschema:"Case ID"`
		ObjectType  string `json:"object_type" jsonschema:"Object type (e.g. cases, assets, ioc, timeline_events, tasks, evidences)"`
		ObjectID    int    `json:"object_id" jsonschema:"Object ID to comment on"`
		CommentText string `json:"comment_text" jsonschema:"Comment text"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_comments_add",
		Description: "Add a comment to a case object",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args commentsAddArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/%s/%d/comments/add", args.ObjectType, args.ObjectID)
		body := map[string]interface{}{"comment_text": args.CommentText}
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Edit comment
	type commentsEditArgs struct {
		CaseID      int    `json:"case_id" jsonschema:"Case ID"`
		ObjectType  string `json:"object_type" jsonschema:"Object type"`
		ObjectID    int    `json:"object_id" jsonschema:"Object ID"`
		CommentID   int    `json:"comment_id" jsonschema:"Comment ID to edit"`
		CommentText string `json:"comment_text" jsonschema:"New comment text"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_comments_edit",
		Description: "Edit an existing comment",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args commentsEditArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/%s/%d/comments/%d/edit", args.ObjectType, args.ObjectID, args.CommentID)
		body := map[string]interface{}{"comment_text": args.CommentText}
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete comment
	type commentsDeleteArgs struct {
		CaseID     int    `json:"case_id" jsonschema:"Case ID"`
		ObjectType string `json:"object_type" jsonschema:"Object type"`
		ObjectID   int    `json:"object_id" jsonschema:"Object ID"`
		CommentID  int    `json:"comment_id" jsonschema:"Comment ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_comments_delete",
		Description: "Delete a comment",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args commentsDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/%s/%d/comments/%d/delete", args.ObjectType, args.ObjectID, args.CommentID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
