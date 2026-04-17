# Numbat Evaluator Skill

## Overview
This MCP server provides a secure, state-aware execution environment for **Numbat**, a statically typed programming language designed specifically for scientific computations, dimensional analysis, and unit conversions. 

Whenever you need to perform mathematical calculations involving physical units, currencies, dates, times, or scientific constants, you should use this MCP server. Numbat's static type system ensures that dimensional equations are mathematically sound, preventing common errors (e.g., adding meters to seconds).

## Tools
The server provides the following tools. Use the discovery tools to look up syntax or capabilities if you are unsure.

* **`numbat_evaluate`**: Evaluates a given Numbat script or expression. It returns the evaluated result, printed standard output, or any compilation/type errors.
* **`numbat_set_variable`**: Defines a persistent variable in the Numbat environment (e.g., saving an intermediate calculation or constant for later use across multiple `numbat_evaluate` calls).
* **`numbat_list_units`**: Returns a comprehensive cheat sheet of all supported physical units, metric/binary prefixes, and live fiat currencies.
* **`numbat_list_functions`**: Returns a reference guide of all built-in mathematical, statistical, date/time, list, and utility functions, including required module imports.
* **`numbat_list_constants`**: Returns a reference guide of all built-in physical (e.g., speed of light, Planck constant) and mathematical (e.g., pi, golden ratio) constants.

---

## Enhanced Function Descriptions & Usage Guidelines

When writing code for the `numbat_evaluate` tool, leverage Numbat's rich standard library and syntax. Below is a guide on how to utilize Numbat's capabilities.

### 1. Self-Discovery & Introspection
If you are asked to perform a calculation but do not know the exact Numbat syntax for a specific unit (e.g., "Is Brazilian Real supported?"), function (e.g., "How do I calculate standard deviation?"), or constant:
* **DO NOT GUESS**: Instead, call `numbat_list_units`, `numbat_list_functions`, or `numbat_list_constants` first.
* Read the returned cheat sheet, then formulate your `numbat_evaluate` call using the correct identifiers.

### 2. Unit Conversions & Formatting
Numbat natively understands physical dimensions (`Length`, `Time`, `Mass`, `Energy`, etc.). Use the `->` (or `to`) operator to perform conversions. 
* **Basic Conversion**: `120 km/h -> mph`
* **Currency**: `4 million BTC -> USD` *(Note: The server fetches live exchange rates on startup)*
* **Temperature**: Temperatures require special care. The base unit is `K` (Kelvin). You can enter and convert using `°C` and `°F`.
  * *Correct*: `25 °C -> °F`
  * *Warning*: Be careful with additions. `10 °C + 1 °C` evaluates to `557.3 K`. Always use `K` for temperature differences: `10 °C + 1 K`.
* **Complex formatting**: Use string interpolation to format outputs clearly.
  ```numbat
  let speed = 2180 km/h
  print("Concorde flew at {speed -> mph:.1f}.")
  ```

### 3. Dimensional Analysis (Typed Holes)
If you know the inputs of a physics equation but forget a constant, you can use Numbat's **Typed Holes** (`?`). The compiler will return an error telling you exactly what physical dimension is missing and suggest available constants.
* **Usage**: `let f: Force = 1 kg * ?` 
* **Server Response**: The server will throw an error indicating `Acceleration` is missing and suggest constants like `g0`. You can use this feedback to correct your formula in a subsequent tool call.

### 4. Date, Time, and Calendar Arithmetic
LLMs often struggle with exact calendar arithmetic and timezones. Delegate this entirely to Numbat.
* **Current Time & Timezones**:
  ```numbat
  now() -> tz("Asia/Tokyo")
  datetime("2024-11-01 12:30:00 Australia/Sydney") -> local
  ```
* **Calendar Math**: Always use `calendar_add` and `calendar_sub` to account for leap years and daylight saving time, rather than just adding `days`.
  ```numbat
  calendar_add(now(), 40 days)
  ```
* **Human Readable Duration**: `1 million seconds -> human` -> Outputs: `"11 days + 13 hours + ..."`

### 5. Advanced Scientific Calculations
Numbat comes with built-in modules for specific scientific domains. 
* **Astronomy (`use extra::astronomy`)**:
  Access properties like `earth_mass`, `solar_radius`, `G`.
  ```numbat
  use extra::astronomy
  let well_depth = G × earth_mass / (g0 × earth_radius) -> km
  ```
* **Chemistry (`use chemistry::elements`)**:
  Query element properties dynamically.
  ```numbat
  use chemistry::elements
  element("Fe").melting_point -> °C
  ```
* **Celestial (`use extra::celestial`)**:
  Compute sunrises, sunsets, and moon phases.
  ```numbat
  use extra::celestial
  sunrise_sunset(Position { lat: 40.713°, lon: -74.006° }, today())
  ```

### 6. Numerical Methods & Algebra
You can solve equations and perform calculus directly using Numbat's numerical methods.
* **Differentiation (`use numerics::diff`)**:
  ```numbat
  use numerics::diff
  fn distance(t) = 0.5 g0 t²
  diff(distance, 2 s, 1e-10 s) # Evaluates velocity at t=2s
  ```
* **Equation Solving (`use numerics::solve`, `use extra::algebra`)**:
  ```numbat
  use extra::algebra
  quadratic_equation(2, -1, -1) # Returns [1, -0.5]
  ```

### 7. Variables, Structs, and Functions
For multi-step reasoning, write a complete Numbat script assigning variables step-by-step.
* **Variables**: `let mass = 50 kg` (Variables are immutable inside a single script; to update, redefine with `let mass = ...`).
* **Functions**: Use `fn` for repeatable logic. Types can usually be inferred, but specifying them guarantees correctness.
  ```numbat
  fn kinetic_energy(mass: Mass, speed: Velocity) -> Energy = 1/2 * mass * speed^2
  ```

### 8. Percentages & Statistics
* **Percentages (`use math::percentage_calculations`)**:
  ```numbat
  100 USD |> increase_by(15%) 
  percentage_change(35 kg, 42 kg)
  ```
* **Statistics (`use math::statistics`)**:
  Use `mean`, `median`, `variance`, `stdev`, `maximum`, `minimum` on arrays of quantities.
  ```numbat
  mean([1 m, 2 m, 300 cm]) # Evaluates to 200 cm
  ```

---

## Best Practices for the Agent

1. **Verify Before Guessing**: Use `numbat_list_functions`, `numbat_list_units`, or `numbat_list_constants` to verify syntax if you are unsure of the correct Numbat identifier.
2. **Always Use Units**: Never pass raw numbers if they represent physical quantities. Instead of `let distance = 5`, use `let distance = 5 meters`. Numbat will automatically handle the underlying conversions.
3. **Use Implicit Multiplication**: Numbat allows implicit multiplication via spaces. `20 kg m / s^2` is perfectly valid and equivalent to `20 * kg * m / s^2`.
4. **Capture Output with `print()`**: If you execute a multi-line script, ensure you `print()` the final result so the MCP server returns it explicitly in the standard output.
5. **Iterate on Type Errors**: If Numbat returns a `Compiler/Runtime Error` (e.g., `Expected Length, found Time`), it means your physical equation is mathematically invalid. Re-evaluate your formula and try calling the tool again.
6. **Use the Pipeline Operator (`|>`)**: For cleaner code, chain function calls using `|>`.
   Example: `17 mph |> round_in(m/s) |> round_in(knots)`

## Example MCP Prompt Injection

**Task**: *Calculate how much energy it takes to heat 1 gallon of water from 70°F to boiling.*

**Tool Call Payload** (`numbat_evaluate`):
```numbat
let density_water = 1 kg / L
let mass_water = 1 gallon × density_water

let c_water = 1 cal / g K
let ΔT = 212 °F - 70 °F

let heat = mass_water × c_water × ΔT

print("Energy to boil 1 gallon of water: {heat -> kJ:.1f}")
print("Or in kWh: {heat -> kWh:.3f}")
```
