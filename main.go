package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/akhenakh/numbat-cgo"
)

// corsMiddleware adds CORS headers so web-based MCP clients (like the MCP Inspector)
// can connect to the server without triggering NetworkErrors.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		transport := serveCmd.String("transport", "stdio", "Transport type: 'stdio', 'sse', or 'streamable'")
		port := serveCmd.Int("port", 8080, "Port for HTTP server (if transport is sse or streamable)")
		serveCmd.Parse(os.Args[2:])

		// Initialize Numbat Engine
		numbatCtx := numbat.NewContext()
		defer numbatCtx.Free() // Ensure Rust memory is freed when server exits

		// Get configured server
		mcpSrv := setupMCPServer(numbatCtx)

		// Always route logs to stderr so JSON-RPC via stdout (stdio) isn't corrupted
		log.SetOutput(os.Stderr)

		if *transport == "streamable" {
			addr := fmt.Sprintf(":%d", *port)
			log.Printf("Starting Numbat MCP Server over Streamable HTTP on %s...", addr)
			log.Printf("Endpoint URL: http://localhost%s/mcp", addr)

			// Create Streamable HTTP Server and wrap it with CORS middleware
			streamableServer := server.NewStreamableHTTPServer(mcpSrv)
			if err := http.ListenAndServe(addr, corsMiddleware(streamableServer)); err != nil {
				log.Fatalf("Streamable HTTP Server error: %v", err)
			}

		} else if *transport == "sse" {
			addr := fmt.Sprintf(":%d", *port)
			log.Printf("Starting Numbat MCP Server over Classic SSE on %s...", addr)
			log.Printf("SSE URL: http://localhost%s/sse", addr)

			// Create Classic SSE Server and wrap it with CORS middleware
			sseServer := server.NewSSEServer(mcpSrv)
			if err := http.ListenAndServe(addr, corsMiddleware(sseServer)); err != nil {
				log.Fatalf("SSE Server error: %v", err)
			}

		} else {
			// Start standard I/O server (Typical for Cursor, Claude Desktop, etc.)
			log.Println("Starting Numbat MCP Server over Stdio...")
			if err := server.ServeStdio(mcpSrv); err != nil {
				log.Fatalf("Stdio Server error: %v", err)
			}
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Numbat MCP Server
Usage: numbat-mcp <command> [options]

Commands:
  serve      Start the MCP Server.
             Usage: numbat-mcp serve [--transport stdio|sse|streamable] [--port 8080]

Examples:
  numbat-mcp serve --transport streamable --port 8080
  numbat-mcp serve --transport sse --port 3000
  numbat-mcp serve --transport stdio
`)
}
