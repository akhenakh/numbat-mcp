package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	numbat "github.com/akhenakh/numbat-cgo"
)

func setupMCPServer(numbatCtx *numbat.Context) *server.MCPServer {
	mcpSrv := server.NewMCPServer("numbat-mcp-server", "1.0.0")

	// TOOL 1: Evaluate Expression
	mcpSrv.AddTool(
		mcp.NewTool("numbat_evaluate",
			mcp.WithDescription(`Executes a script or expression in Numbat, a statically typed language for scientific computations, dimensional analysis, and unit conversions.

WHEN TO USE:
Use this for ALL math, physics, currency, date/time, and dimensional calculations to prevent hallucination and guarantee mathematical correctness.

HOW TO USE:
1. Units are mandatory: Append units to numbers (e.g., '50 km / 2 hours').
2. Conversions: Use the '->' operator (e.g., '120 km/h -> mph', '100 EUR -> USD').
3. Standard Library: You can use imports for advanced functions:
   - 'use extra::astronomy' (for G, earth_mass, solar_radius)
   - 'use chemistry::elements' (for element("Fe").melting_point)
   - 'use numerics::diff', 'use numerics::solve', 'use extra::algebra'
4. Dates & Time: Use 'now() -> tz("Asia/Tokyo")' or 'calendar_add(now(), 40 days)'.
5. Typed Holes: Put '?' in an equation if you forget a constant (e.g., 'let f: Force = 1 kg * ?'). The compiler will error and suggest the right constant.
6. Temperature: Additions require Kelvin (e.g., '10 °C + 1 K'). Conversions are direct (e.g., '25 °C -> °F').`),
			mcp.WithString("expression",
				mcp.Required(),
				mcp.Description("The raw Numbat code to execute. Can be a single expression (e.g., '120 km/h -> m/s') or a sequence of statements (e.g., 'let d = 10 m; d / 2 s'). Do not wrap in markdown code blocks."),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			expr := request.GetString("expression", "")

			res, err := numbatCtx.Interpret(expr)
			if err != nil {
				// Return the error as a tool result so the agent knows what went wrong and can iterate/fix it
				return mcp.NewToolResultError(fmt.Sprintf("Numbat Compiler/Runtime Error:\n%v", err)), nil
			}

			// Format the response with detailed, structured output for better LLM consumption
			var output string
			if res.IsQuantity {
				output = fmt.Sprintf("Result: %s\n\nNumeric Value: %f\nUnit: %s", res.StringOutput, res.Value, res.Unit)
			} else {
				output = fmt.Sprintf("Result: %s", res.StringOutput)
			}

			return mcp.NewToolResultText(output), nil
		},
	)

	// TOOL 2: Set Variable
	mcpSrv.AddTool(
		mcp.NewTool("numbat_set_variable",
			mcp.WithDescription(`Define a persistent variable in the Numbat environment.
Use this to safely store intermediate results, physical properties, or constants that will be reused in subsequent 'numbat_evaluate' calls. 
(Note: You can also define temporary variables directly inside 'numbat_evaluate' using the 'let' keyword).`),
			mcp.WithString("name", mcp.Required(), mcp.Description("The exact name of the variable (e.g., 'distance_to_mars').")),
			mcp.WithNumber("value", mcp.Required(), mcp.Description("The numerical value of the variable (e.g., 225000000).")),
			mcp.WithString("unit", mcp.Description("The physical unit of the variable (e.g., 'km', 'm/s', 'kg'). Leave empty if dimensionless.")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			value := request.GetFloat("value", 0.0)
			unit := request.GetString("unit", "") // defaults to empty string if missing

			err := numbatCtx.SetVariable(name, value, unit)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to set variable: %v", err)), nil
			}

			unitDisplay := unit
			if unitDisplay == "" {
				unitDisplay = "(dimensionless)"
			}

			msg := fmt.Sprintf("Variable '%s' set successfully to %f %s.", name, value, unitDisplay)
			return mcp.NewToolResultText(msg), nil
		},
	)

	return mcpSrv
}
