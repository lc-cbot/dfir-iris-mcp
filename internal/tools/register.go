package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"dfir-iris-mcp/internal/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterAll(s *mcp.Server, c *client.Client) {
	registerSystem(s, c)
	registerSettings(s, c)
	registerCases(s, c)
	registerAlerts(s, c)
	registerAssets(s, c)
	registerNotes(s, c)
	registerIOCs(s, c)
	registerTimeline(s, c)
	registerTasks(s, c)
	registerEvidences(s, c)
	registerDatastore(s, c)
	registerComments(s, c)
	registerUsers(s, c)
	registerGroups(s, c)
	registerCustomers(s, c)
}

func textResult(data json.RawMessage) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}
}

func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: err.Error()}},
		IsError: true,
	}
}

func cidQuery(caseID int) map[string]string {
	return map[string]string{"cid": strconv.Itoa(caseID)}
}

// toBody converts a struct to a map for use as a JSON request body,
// excluding specified keys and nil values.
func toBody(args interface{}, exclude ...string) map[string]interface{} {
	b, _ := json.Marshal(args)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	ex := make(map[string]bool)
	for _, k := range exclude {
		ex[k] = true
	}
	for k, v := range m {
		if ex[k] || v == nil {
			delete(m, k)
		}
	}
	return m
}

// toQuery converts a struct to a query-string map,
// excluding specified keys and nil values.
func toQuery(args interface{}, exclude ...string) map[string]string {
	b, _ := json.Marshal(args)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	ex := make(map[string]bool)
	for _, k := range exclude {
		ex[k] = true
	}
	q := make(map[string]string)
	for k, v := range m {
		if !ex[k] && v != nil {
			q[k] = fmt.Sprint(v)
		}
	}
	return q
}
