# Numbat MCP Server 🧮🦇

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that provides AI agents and LLMs with a mathematically rigorous, unit-aware computation engine powered by [Numbat](https://numbat.dev/).

Large Language Models are notoriously unreliable at math, complex unit conversions, and date/time arithmetic. This server solves that problem by offloading these tasks to Numbat, a statically typed programming language designed specifically for scientific computations, dimensional analysis, and unit-safe math.

## Features

* **Unit-Safe Computations**: Add, subtract, multiply, and divide quantities safely. (e.g., `120 km/h -> mph`).
* **Live Currency Conversions**: Perform fiat and crypto conversions (`100 EUR -> USD`).
* **Date & Time Arithmetic**: Native timezone conversions and calendar math (`calendar_add(now(), 40 days)`).
* **Extensive Scientific Library**: Built-in constants and functions for astronomy, chemistry, statistics, numerical methods, and more.
* **LLM Self-Discovery**: Includes tools for the LLM to dynamically list available units, functions, and constants so it never has to guess the syntax.
* **Persistent State**: Set and retrieve variables across multiple tool calls within a chat session.

## Installation

Download the binary for your operating system, and configure your agentic system, see below.

### Building it

Ensure you have [Go](https://go.dev/doc/install) installed, along with the required C/C++ toolchains for building the underlying Rust/C-bindings if necessary.

```bash
# Clone the repository
git clone https://github.com/akhenakh/numbat-mcp.git
cd numbat-mcp

# Build the binary
go build -o numbat-mcp .

mv numbat-mcp /usr/local/bin/
```

## Usage

The server supports three MCP transports: `stdio` (standard for most local clients), `sse`, and `streamable`.

```bash
# Start via stdio (Default for Claude Desktop, Cursor, Opencode)
numbat-mcp serve --transport stdio

# Start via SSE on a specific port
numbat-mcp serve --transport sse --port 3000

# Start via Streamable HTTP
numbat-mcp serve --transport streamable --port 8080
```

## Client Configuration

### Opencode
To enable the Numbat MCP server in Opencode, add the following configuration to your MCP settings section:

```yaml
    "numbat-mcp": {
      "type": "local",
      "command": ["numbat-mcp", "serve", "--transport", "stdio"],
      "enabled": true,
    },
```
*(Make sure the `numbat-mcp` binary is in your system's `PATH`, or provide the absolute path in the `command` array).*

### Claude Desktop
To use this with the Claude Desktop app, add the following to your `claude_desktop_config.json` (located at `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "numbat": {
      "command": "numbat-mcp",
      "args": ["serve", "--transport", "stdio"]
    }
  }
}
```

## Available MCP Tools

Once connected, your AI agent will have access to the following tools:

1. **`numbat_evaluate`**: Evaluates a Numbat script or expression. Supports full dimensional analysis and typed holes (`?`).
2. **`numbat_set_variable`**: Defines a persistent variable in the Numbat environment to be used in future evaluations.
3. **`numbat_list_units`**: Returns a cheat sheet of all supported physical units, currencies, and prefixes.
4. **`numbat_list_functions`**: Returns a reference guide of all built-in mathematical, date/time, and utility functions.
5. **`numbat_list_constants`**: Returns a list of all built-in physical and mathematical constants.

## Example Prompts for your LLM

Once the server is connected, you can ask your LLM things like:
* *"If I drive at 80 mph for 3.5 hours, how many kilometers did I travel?"*
* *"How much energy does it take to boil 1 gallon of water from 70°F?"*
* *"Convert 50,000 Japanese Yen to Euros."*
* *"What is the gravitational pull of Jupiter at its surface?"*

The LLM will automatically format these into Numbat syntax, execute them on the server, and return the mathematically proven result to you.

---
**Credits**: Powered by [Numbat](https://github.com/sharkdp/numbat) (by David Peter) and [mcp-go](https://github.com/mark3labs/mcp-go) (by Mark3Labs).
