package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerTasks(s *mcp.Server, c *client.Client) {
	// List tasks
	type tasksListArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_tasks_list",
		Description: "List all tasks in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tasksListArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/case/tasks/list", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get task
	type tasksGetArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		TaskID int `json:"task_id" jsonschema:"Task ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_tasks_get",
		Description: "Get details of a specific task in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tasksGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/tasks/%d", args.TaskID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add task
	type tasksAddArgs struct {
		CaseID          int     `json:"case_id" jsonschema:"Case ID"`
		TaskTitle       string  `json:"task_title" jsonschema:"Title of the task"`
		TaskDescription *string `json:"task_description,omitempty" jsonschema:"Task description"`
		TaskAssigneesID *[]int  `json:"task_assignees_id,omitempty" jsonschema:"List of user IDs to assign"`
		TaskStatusID    *int    `json:"task_status_id,omitempty" jsonschema:"Task status ID"`
		TaskTags        *string `json:"task_tags,omitempty" jsonschema:"Comma-separated tags"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_tasks_add",
		Description: "Add a new task to a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tasksAddArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Post(ctx, "/case/tasks/add", cidQuery(args.CaseID), toBody(args, "case_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update task
	type tasksUpdateArgs struct {
		CaseID          int     `json:"case_id" jsonschema:"Case ID"`
		TaskID          int     `json:"task_id" jsonschema:"Task ID to update"`
		TaskTitle       *string `json:"task_title,omitempty" jsonschema:"New task title"`
		TaskDescription *string `json:"task_description,omitempty" jsonschema:"New description"`
		TaskAssigneesID *[]int  `json:"task_assignees_id,omitempty" jsonschema:"New list of assignee user IDs"`
		TaskStatusID    *int    `json:"task_status_id,omitempty" jsonschema:"New status ID"`
		TaskTags        *string `json:"task_tags,omitempty" jsonschema:"New comma-separated tags"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_tasks_update",
		Description: "Update a task in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tasksUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/tasks/update/%d", args.TaskID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "task_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete task
	type tasksDeleteArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		TaskID int `json:"task_id" jsonschema:"Task ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_tasks_delete",
		Description: "Delete a task from a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args tasksDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/tasks/delete/%d", args.TaskID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
