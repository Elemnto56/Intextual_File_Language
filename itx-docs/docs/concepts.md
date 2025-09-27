## ISEC
ISEC stands for "Intext's Script Execution Core". It's Intext's execution engine that takes your code, and outputs it properly. It's chained by:
1. Getting your source code
<pre>output "Hello!";</pre>

2. Lexing/tokening it into tokens for the Parser to read
```JSON
[
  {
    "LINE": 1,
    "TYPE": "KEYWORD",
    "VAL": "output"
  },
  {
    "LINE": 1,
    "TYPE": "STRING",
    "VAL": "Hello!"
  },
  {
    "LINE": 1,
    "TYPE": "SYMBOL",
    "VAL": ";"
  }
]
```

3. Then the Parser parses it into an Abstract Syntax Tree *(AST)*, whilst checking for correct syntax
```JSON
[
  {
    "line": 1,
    "meta": {
      "print_type": "simple",
      "raw_type": "STRING"
    },
    "type": "output",
    "value": "Hello!"
  }
]
```

4. Then running it through a phase called a "Validator" which makes sure called variables were declared and function args are correct

5. The Interpreter then executes it based on what's given
```
Hello!
```

## Bound System
***Note: This section is subject to change***

The Bound System is a system that checks the relativity of your code. Relative to whether if it's only used in the script itself, or if it interacts with the outside world. There are three categories that code could align with.

### In-Bounds 
In-Bounds, means all code or funcs are only within the script. They do that interact with the outside, only the within the script itself. Some examples include:
- `output`
- `let`
- *logic based syntax such as if statements and loops*

### Passive Bounds
Passive Bounds are code or funcs that not only interact with the script, but they also interact with the outside world. Examples include: <!--Add input when possible -->

- `read()`

To clear any confusion, `read()` is considered one because it interacts with the outside by getting the contents of file, whilst assigning those contents to a value that will be given to variable.

### Out-Bounds
Out-Bounds are code or funcs that ***do not*** interact with the script whatsoever. They interact with outside files or with data that doesn't involve Intext. Examples include:
- `write()`
- `append()`

To clear any confusion, yes, `write()` and `append()` do need variables for input which could consider them to be In-Bounds, but, their purpose is enacted outside of the script.