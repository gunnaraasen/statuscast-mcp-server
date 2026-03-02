package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/gunnaraasen/statuscast-mcp-server/internal/config"
	"github.com/gunnaraasen/statuscast-mcp-server/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	c := client.New(cfg.Domain, cfg.Token)

	s := mcp.NewServer(&mcp.Implementation{
		Name:    "statuscast-mcp-server",
		Version: "1.0.0",
	}, nil)

	s.AddReceivingMiddleware(func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (result mcp.Result, err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			return next(ctx, method, req)
		}
	})

	tools.RegisterAll(s, c)

	switch cfg.Transport {
	case "http":
		handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server { return s }, nil)
		log.Printf("Listening on :%s", cfg.Port)
		log.Fatal(http.ListenAndServe(":"+cfg.Port, handler))
	default: // "stdio"
		if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Printf("server error: %v", err)
		}
	}
}
