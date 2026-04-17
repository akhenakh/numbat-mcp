package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	numbat "github.com/akhenakh/numbat-cgo"
)

func setupMCPServer(numbatCtx *numbat.Context) *server.MCPServer {
	mcpSrv := server.NewMCPServer("numbat-mcp-server", "1.1.0")

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
				return mcp.NewToolResultError(fmt.Sprintf("Numbat Compiler/Runtime Error:\n%v", err)), nil
			}

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
Use this to safely store intermediate results, physical properties, or constants that will be reused in subsequent 'numbat_evaluate' calls.`),
			mcp.WithString("name", mcp.Required(), mcp.Description("The exact name of the variable (e.g., 'distance_to_mars').")),
			mcp.WithNumber("value", mcp.Required(), mcp.Description("The numerical value of the variable (e.g., 225000000).")),
			mcp.WithString("unit", mcp.Description("The physical unit of the variable (e.g., 'km', 'm/s', 'kg'). Leave empty if dimensionless.")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := request.GetString("name", "")
			value := request.GetFloat("value", 0.0)
			unit := request.GetString("unit", "")

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

	// TOOL 3: List Units
	mcpSrv.AddTool(
		mcp.NewTool("numbat_list_units",
			mcp.WithDescription("Returns a comprehensive list of all supported physical units, currencies, and their aliases in Numbat."),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText(numbatUnitsDoc), nil
		},
	)

	// TOOL 4: List Functions
	mcpSrv.AddTool(
		mcp.NewTool("numbat_list_functions",
			mcp.WithDescription("Returns a list of all built-in mathematical, statistical, date/time, and utility functions available in Numbat."),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText(numbatFunctionsDoc), nil
		},
	)

	// TOOL 5: List Constants
	mcpSrv.AddTool(
		mcp.NewTool("numbat_list_constants",
			mcp.WithDescription("Returns a list of all built-in physical and mathematical constants available in Numbat."),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultText(numbatConstantsDoc), nil
		},
	)

	return mcpSrv
}

// -----------------------------------------------------------------------------
// Reference Documentation Constants
// -----------------------------------------------------------------------------

const numbatUnitsDoc = `# Supported Units in Numbat

Numbat supports SI metric prefixes (milli, centi, kilo, mega, etc.) and binary prefixes (kibi, mebi, etc.) for applicable units.

## Base & Common Physical Units
- **Length**: meter (m), centimeter (cm), kilometer (km), inch (in), foot (ft), yard (yd), mile (mi), nautical_mile (nmi), lightyear (ly), parsec (pc), angstrom (Å), astronomicalunit (au), fathom.
- **Area**: are, hectare (ha), acre, barn.
- **Volume**: liter (L), gallon (gal), cup, fluidounce (floz), pint, quart, tablespoon (tbsp), teaspoon (tsp), barrel.
- **Mass**: gram (g), kilogram (kg), tonne, pound (lb), ounce (oz), stone, dalton (Da), grain.
- **Time**: second (s), minute (min), hour (h), day (d), week, month, year, decade, century, millennium.
- **Temperature**: kelvin (K), celsius (°C), fahrenheit (°F). *(Note: °C and °F are for absolute temperatures/conversions. Use K for math).*
- **Current / Voltage / Resistance**: ampere (A), volt (V), ohm (Ω).
- **Force / Pressure**: newton (N), pascal (Pa), bar, atmosphere (atm), psi, mmHg, torr, dyne.
- **Energy / Power**: joule (J), electronvolt (eV), calorie (cal), BTU, watt (W), horsepower (hp).
- **Frequency / Angle**: hertz (Hz), rpm, degree (°), radian (rad), arcmin (′), arcsec (″).
- **Information**: bit, byte (B), KB, MB, etc.

## Currencies (Live Exchange Rates)
- **USD** ($), **EUR** (€), **GBP** (£), **JPY** (¥), **AUD**, **CAD**, **CHF**, **CNY** (元), **INR** (₹), **NZD**, etc.

## Astronomy (requires 'use extra::astronomy')
- **Length**: earth_radius, solar_radius, lunar_radius, jupiter_radius.
- **Mass**: earth_mass, solar_mass, lunar_mass, jupiter_mass.
`

const numbatFunctionsDoc = `# Built-in Functions in Numbat

## Mathematics & Numerics
- **Basic**: ` + "`abs(x)`" + `, ` + "`sqrt(x)`" + `, ` + "`cbrt(x)`" + `, ` + "`sqr(x)`" + `, ` + "`round(x)`" + `, ` + "`floor(x)`" + `, ` + "`ceil(x)`" + `, ` + "`mod(a, b)`" + `
- **Rounding to Unit**: ` + "`round_in(m, 5.3 m)`" + `, ` + "`floor_in(cm, ...)`" + `
- **Transcendental**: ` + "`exp(x)`" + `, ` + "`ln(x)`" + `, ` + "`log10(x)`" + `, ` + "`log2(x)`" + `
- **Trigonometry**: ` + "`sin(x)`" + `, ` + "`cos(x)`" + `, ` + "`tan(x)`" + `, ` + "`asin(x)`" + `, ` + "`acos(x)`" + `, ` + "`atan2(y, x)`" + `
- **Calculus & Solving** *(requires imports)*: 
  - ` + "`use numerics::diff; diff(f, x, dx)`" + `
  - ` + "`use numerics::solve; root_bisect(...)`, `root_newton(...)`" + `
  - ` + "`use extra::algebra; quadratic_equation(a, b, c)`" + `

## Statistics & Combinatorics
- **Stats**: ` + "`minimum(xs)`" + `, ` + "`maximum(xs)`" + `, ` + "`mean(xs)`" + `, ` + "`median(xs)`" + `, ` + "`variance(xs)`" + `, ` + "`stdev(xs)`" + `
- **Combinatorics**: ` + "`factorial(n)`" + ` (or ` + "`n!`" + `), ` + "`binom(n, k)`" + `, ` + "`fibonacci(n)`" + `
- **Random**: ` + "`random()`" + `, ` + "`rand_uniform(a, b)`" + `, ` + "`rand_int(a, b)`" + `, ` + "`rand_norm(μ, σ)`" + `

## Date & Time
- **Current**: ` + "`now()`" + `, ` + "`today()`" + `
- **Parsing**: ` + "`datetime(\"2024-01-01 12:00 UTC\")`" + `, ` + "`date(\"2024-01-01\")`" + `
- **Timezones**: ` + "`now() -> tz(\"Asia/Tokyo\")`" + `, ` + "`datetime(...) -> local`" + `
- **Calendar Math**: ` + "`calendar_add(now(), 2 years)`" + `, ` + "`calendar_sub(now(), 3 days)`" + `
- **Formatting**: ` + "`1 million seconds -> human`" + ` (Outputs: X days + Y hours...)
- **Unix Epoch**: ` + "`now() -> unixtime_s`" + `, ` + "`from_unixtime_s(1658346725)`" + `

## Lists & Strings
- **Lists**: ` + "`len(xs)`" + `, ` + "`head(xs)`" + `, ` + "`tail(xs)`" + `, ` + "`map(f, xs)`" + `, ` + "`filter(f, xs)`" + `, ` + "`sum(xs)`" + `, ` + "`join(xs, sep)`" + `
- **Strings**: ` + "`str_length(s)`" + `, ` + "`lowercase(s)`" + `, ` + "`uppercase(s)`" + `, ` + "`str_replace(pat, rep, s)`" + `
- **Formatting Conversions**: ` + "`42 -> hex`" + `, ` + "`42 -> bin`" + `, ` + "`42 -> oct`" + `
`

const numbatConstantsDoc = `# Built-in Constants in Numbat

## Mathematical
- **pi / π**: 3.14159...
- **tau / τ**: 6.28318...
- **e**: 2.71828...
- **golden_ratio / φ**: 1.61803...

## Named Numbers
- ` + "`hundred`" + `, ` + "`thousand`" + `, ` + "`million`" + `, ` + "`billion`" + `, ` + "`trillion`" + `, ` + "`quadrillion`" + `
- ` + "`half`" + `, ` + "`quarter`" + `, ` + "`double`" + `, ` + "`dozen`" + `

## Physics Constants
- **speed_of_light / c**: Speed of light in vacuum.
- **gravitational_constant / G**: Newtonian constant of gravitation.
- **gravity / g0**: Standard acceleration of gravity on earth.
- **planck_constant / ℎ**: The Planck constant.
- **h_bar / ℏ**: The reduced Planck constant.
- **electron_mass**: Mass of the electron.
- **proton_mass**: Mass of the proton.
- **neutron_mass**: Mass of the neutron.
- **elementary_charge / electron_charge**: Charge of the electron.
- **magnetic_constant / µ0**: Vacuum magnetic permeability.
- **electric_constant / ε0**: Vacuum electric permittivity.
- **avogadro_constant / N_A**: Avogadro's number.
- **boltzmann_constant / k_B**: Boltzmann constant.
- **gas_constant / R**: Ideal gas constant.
- **stefan_boltzmann_constant**: Stefan-Boltzmann constant.
- **fine_structure_constant / α**: Fine structure constant.

## Scales / Planck Units
- ` + "`planck_length`" + `, ` + "`planck_mass`" + `, ` + "`planck_time`" + `, ` + "`planck_temperature`" + `, ` + "`planck_energy`" + `
- ` + "`bohr_radius`" + `, ` + "`rydberg_constant`" + `
`
