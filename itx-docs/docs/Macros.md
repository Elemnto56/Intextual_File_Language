# Macros

## Overview

Macros are reusable functions that provide a powerful mechanism for code reusability and shareability. They allow you to define custom functionality that can be invoked throughout your script, promoting cleaner code organization and reducing duplication.


***Note***: Some versions may have these features while others do not. Check the changelog to see which versions do or don't have them. If the Changelog isn't descriptive enough for you, well, I'm sorry, I'm doing my best ¯\\_(ツ)_/¯

<hr></hr>

## Syntax

The basic syntax for defining a macro follows this pattern:

```
macro (callParameter: type) . macroName(parameters) -> returnType {
    // macro body
}
```

**Components:**

- `macro` - The keyword that declares a macro
- `(callParameter: type)` - The call parameter that determines how the macro is invoked
- `.macroName(parameters)` - The macro name and any additional parameters
- `-> returnType` - Optional return type specification
- `{ }` - The macro body containing the implementation
<hr></hr>

## Examples

### Basic Macro with Output

```
// Define a simple macro that outputs text
let x: string = "Hello, World!"

macro (call: string) . mac() {
    output call
}

// Call the macro
x.mac()
```

### Macro with Parameters

```
// Define a macro that accepts parameters
macro . mac(num: int) {
    output num
}

// Call with parameter
mac(2)
```

**Note:** The period "." is part of the syntax regardless of whether you're using call parameters.

### Macro with Return Value

```
// Define a macro that returns a value
macro . mac() -> string {
    return "What"
}

// Use the returned value
x = mac()
```

<hr></hr>

## Usage Patterns

### Call Parameter Macros
When you define a macro with a call parameter, you invoke it by calling the method on a variable of the matching type:

```
let myString: string = "Hello"
macro (call: string) . greet() {
    output call
}
myString.greet()  // Outputs: "Hello"
```

### Standalone Macros
Macros without call parameters can be invoked directly:

```
macro . sayHello() {
    output "Hello, World!"
}
sayHello()  // Direct invocation
```

## Important Notes

- **Type matching**: Call parameters must match the type of the variable used to invoke the macro
- **Flexibility**: Macros can be defined with or without parameters, with or without return values, and with or without call parameters, providing maximum flexibility for different use cases

## Best Practices

- Use descriptive macro names that clearly indicate their purpose
- Use standalone macros for utility functions that don't require a specific calling context
- Document (using comments) complex macros to ensure they remain maintainable and shareable