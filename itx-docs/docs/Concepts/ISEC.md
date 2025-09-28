# ISEC
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