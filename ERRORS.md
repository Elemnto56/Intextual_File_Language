## Intext Error Reference

A list of all possible errors that can occur during Intext execution and validation.

---

### Syntax Errors

| Error Name        | Description                                               | Example        | Suggested Fix                      |
|-------------------|-----------------------------------------------------------|----------------|------------------------------------|
| MissingBreaker  | A statement was not properly terminated with a semicolon or break symbol. | output "Hello" | Add a semicolon: `output "Hello";` |
| LexerError | The Lexer encountered an error which usually involves an invalid character being detected. | decl@re x: int = 5; | Remove the character and replace it with something valid. |
| MalformedSyntax | The following statement had bad syntax. Rules were broken. | `declare x: int 5` | Double check your work, and fix mistakes. |

---

### Validation Errors

| Error Name           | Description                                          | Example                  | Suggested Fix                 |
|----------------------|------------------------------------------------------|--------------------------|-------------------------------|
| TypeMismatch         | A variable was assigned a value of a different type. | declare: int x = "hello"; | Assign a matching type value. |
| IllegalMutation | Tried to assign something to either a literal or immuatable variable. | "Hell" += "o" | Make sure to check both sides of the assigning and know that they are both mutable. |
| VariableNotFound | A variable that did not exist was called. | `output myVar;` | Make sure the variable you are calling is declared. |
| UnknownValue | A value assigned to a variable was so bad, it couldn't be TypeMismatched. | `declare x: string = hello;` | Make sure to check your syntax. |
---

### Runtime Errors

| Error Name       | Description                                        | Example            | Suggested Fix                     |
|------------------|----------------------------------------------------|--------------------|-----------------------------------|
| FileError | The following file could either not be found, opened, or written to. | N/A | Chances are, you are on Linux, thus, chmod it, and make it `666` or fix the file path. |
| RangeException | Tried to access either an index too high or low on an order | ```declare x: order = ["apples", "oranges"];   output x[10];``` | Try to fix any loops or conditions that would lead to an order having an index too low or high.|
| DivisionByZero | You tried to divide by zero. If you know math, then you know this is bad. | `output 0/0;` | If it was by a variable, then try to refactor your code to not make it zero |

---

### Other

| Error Name      | Description                            | Example       | Suggested Fix         |
|-----------------|----------------------------------------|---------------|-----------------------|
<!-- Add more here when needed>