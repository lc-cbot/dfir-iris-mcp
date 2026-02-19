package tools

import (
	"context"
	"fmt"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerNotes(s *mcp.Server, c *client.Client) {
	// List note directories (note groups)
	type notesDirsListArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_groups_list",
		Description: "List all note directories (groups) in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesDirsListArgs) (*mcp.CallToolResult, any, error) {
		data, err := c.Get(ctx, "/case/notes/directories/filter", cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add note directory
	type notesDirsAddArgs struct {
		CaseID int    `json:"case_id" jsonschema:"Case ID"`
		Name   string `json:"name" jsonschema:"Name of the note directory"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_groups_add",
		Description: "Create a new note directory (group) in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesDirsAddArgs) (*mcp.CallToolResult, any, error) {
		body := map[string]interface{}{"name": args.Name}
		data, err := c.Post(ctx, "/case/notes/directories/add", cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update note directory
	type notesDirsUpdateArgs struct {
		CaseID      int    `json:"case_id" jsonschema:"Case ID"`
		DirectoryID int    `json:"directory_id" jsonschema:"Note directory ID to update"`
		Name        string `json:"name" jsonschema:"New directory name"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_groups_update",
		Description: "Update a note directory (group) in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesDirsUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/notes/directories/update/%d", args.DirectoryID)
		body := map[string]interface{}{"name": args.Name}
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete note directory
	type notesDirsDeleteArgs struct {
		CaseID      int `json:"case_id" jsonschema:"Case ID"`
		DirectoryID int `json:"directory_id" jsonschema:"Note directory ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_groups_delete",
		Description: "Delete a note directory from a case (deletes all notes in it)",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesDirsDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/notes/directories/delete/%d", args.DirectoryID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Get note
	type notesGetArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		NoteID int `json:"note_id" jsonschema:"Note ID"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_get",
		Description: "Get details of a specific note",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesGetArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/notes/%d", args.NoteID)
		data, err := c.Get(ctx, path, cidQuery(args.CaseID))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Add note
	type notesAddArgs struct {
		CaseID      int    `json:"case_id" jsonschema:"Case ID"`
		NoteTitle   string `json:"note_title" jsonschema:"Title of the note"`
		NoteContent string `json:"note_content" jsonschema:"Content of the note (supports markdown)"`
		DirectoryID int    `json:"directory_id" jsonschema:"Note directory ID to add the note to"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_add",
		Description: "Add a new note to a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesAddArgs) (*mcp.CallToolResult, any, error) {
		body := map[string]interface{}{
			"note_title":   args.NoteTitle,
			"note_content": args.NoteContent,
			"directory_id": args.DirectoryID,
		}
		data, err := c.Post(ctx, "/case/notes/add", cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Update note
	type notesUpdateArgs struct {
		CaseID      int     `json:"case_id" jsonschema:"Case ID"`
		NoteID      int     `json:"note_id" jsonschema:"Note ID to update"`
		NoteTitle   *string `json:"note_title,omitempty" jsonschema:"New note title"`
		NoteContent *string `json:"note_content,omitempty" jsonschema:"New note content"`
		DirectoryID *int    `json:"directory_id,omitempty" jsonschema:"Move note to a different directory"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_update",
		Description: "Update an existing note in a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesUpdateArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/notes/update/%d", args.NoteID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), toBody(args, "case_id", "note_id"))
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Delete note
	type notesDeleteArgs struct {
		CaseID int `json:"case_id" jsonschema:"Case ID"`
		NoteID int `json:"note_id" jsonschema:"Note ID to delete"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_delete",
		Description: "Delete a note from a case",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesDeleteArgs) (*mcp.CallToolResult, any, error) {
		path := fmt.Sprintf("/case/notes/delete/%d", args.NoteID)
		data, err := c.Post(ctx, path, cidQuery(args.CaseID), nil)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})

	// Search notes
	type notesSearchArgs struct {
		CaseID     int    `json:"case_id" jsonschema:"Case ID"`
		SearchTerm string `json:"search_term" jsonschema:"Text to search for in notes"`
	}
	mcp.AddTool(s, &mcp.Tool{
		Name:        "dfir_iris_notes_search",
		Description: "Search notes in a case by keyword",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args notesSearchArgs) (*mcp.CallToolResult, any, error) {
		body := map[string]interface{}{"search_term": args.SearchTerm}
		data, err := c.Post(ctx, "/case/notes/search", cidQuery(args.CaseID), body)
		if err != nil {
			return errorResult(err), nil, nil
		}
		return textResult(data), nil, nil
	})
}
