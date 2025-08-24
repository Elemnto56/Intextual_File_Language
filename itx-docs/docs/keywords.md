## output
`output` allows the user to print data to the console. `output` has multiple methods of itself.

- ### Simple Output
    This method allows one-valued pieces of data to be printed. Wheter it's a variable, `string`, `int`, or `bool`. As long as it's a single value, it's considered Simple Output.
    <br>
    <b>Example:</b>
    <pre class="cd">output "Hello World";</pre>

- ### Mixed Output
    Sometimes will be referred as "Spaghetti Output" involves multiple values in ``output``. Mixed Output is very freeform when it comes to what can or cannot be printed. However, you must use commas to seperate each value. Using something like a plus (`+`) would make ISEC *([What's ISEC?](concepts.md))* think you are doing math, which would then cause an error.
    <br>
    You are able to do math in `output`, but mixing math and Mixed Output together will cause an error.
    
    <b>Error example:</b>
    <pre>output 2 + 2, "Is equal to four";</pre>
    <b>Correct Example:</b>
    <pre>output "I love dogs", "and", "cats";</pre>

- ### Mathematical Output
    This method's only purpose is to solve math. You can use variables in it, but not literals. Failing to do so, will lead to an error.
    <b>Error Example:</b>
    <pre>output 2 + "Hello!";</pre>
    <b>Correct Example:</b>
    <pre>output 2 + 5 / 10 * x; <span style="color: grey;">// x is a variable in this case</span></pre>


## let, declare
`let` allows the user to create variables. Variables can be used for a multitude of things. Also, to honour the Legacy version of Intext, besides using `let` to declare a variable, you can use `declare`. Besides what keyword you can use, similar to `output`, `let` has methods.

- ### Concat
    Allows multiple values to be in one variable. To do so, seperate each value by a comma <br>--> <code>"Hello", " World"</code>
    <br>
    <br>
    However, one could also use plus symbols, BUT only if all of the other values are of type `string`. Reason why, is because ISEC expects values that have `+` seperating them, to be only of a type `int`. Since, `+` is usually used in mathematics. However, when it's of type `string`, ISEC concats. 
    <br>
    <br>
    Now, one may wonder, where does the error occur? Well with concat, and even Mixed Output, this is legal.
    <pre>let x: string = "Hello", true, 5, "World", '!';</pre> 
    As jumbled as that mess is, ISEC will evaluate it. In fact, it helps when reading and writing files, to prevent conversion errors. But this however, is not illegal.
    <pre>let x: string = "Hello" + 5;</pre>
    ISEC thinks you are trying to do math since you added an `int`. So when it tries to add, it doesn't work. To fix it, either make 5 a `string` or use commas.

    <b>Example:</b>
    <pre>let x: string = "The world population is", 8.1, 'B';</pre>

- ### Mathematics
    Gives the value of the math you input to the variable you want. Suppose you do
    <pre>let x: int = 4 + 2;</pre>
    Then `x` will equal 6.

    You can even incorporate other variables in the math, so something like
    <pre>
    let x: int = 4;
    let y: int = 10 + x;
    </pre>
    Would have `y` equal 14.

- ### Functions
    You can use `let` to assign the value of functions to a variable. For example, with `read()` if you do 
    <pre>let x: string = read("file.txt")</pre>
    Then `x` will equal the contents of file.txt. **HOWEVER**, there are caveats. Whenever you use `read()`, or basically any function, for safety, you can only have it be assigned to a variable. It's value may not be directly outputed, or used in a concat. For example, doing 
    <pre>output read("file.txt")</pre>
    Is illegal, and would result in an error. What to do instead, is to assgin read to a variable, and then output. 
    
    This is all done so ISEC can properly assess each variable, preventing huge critical errors in scripts. So for example, say you `output` a concat of `read()` and some other strings (ex: <pre>output read("file.txt"), "is my friend", '!')</pre> But then it turns out that file.txt doesn't exist! How would ISEC deal with that, without affecting output? So, this is why you must put functions in a variable first before using their value.