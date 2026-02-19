package main

import (
	"context"
	"log"

	"dfir-iris-mcp/internal/client"
	"dfir-iris-mcp/internal/config"
	"dfir-iris-mcp/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	c := client.New(cfg.BaseURL, cfg.APIKey)

	s := mcp.NewServer(
		&mcp.Implementation{Name: "dfir-iris-mcp", Version: "1.0.0"},
		nil,
	)

	tools.RegisterAll(s, c)

	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server: %v", err)
	}
}
