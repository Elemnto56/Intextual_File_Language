This page will go into detail about *most* functions that Intext has built-in

### read()

- **Purpose**
    Used to read file content

- **Syntax**
    
    ```
    read(filename)
    ```
    
    ***Note:*** Event thought the example shows it being used in a standalone context, ``read()`` doesn't work this way. You must assign it's value to a variable. The only reason why the syntax is show like this, is to focus on mainly `read()`.

    - ***filename***
        *filename* must be a string. Whether variable or a literal, it must be a string, and the name of a real file. The *filename* can also be a path to one (i.e. `"../myFile.txt"`). An error will be thrown if the file doesn't exist, or if the type is not a string.

- **Caveats**
    - ``read()`` has the ability to read files outside of UTF-8, meaning it can read binary. So, suppose you decide to read an executable from a C++ compilation. ``read()`` will spam your terminal with junk characters which, in some cases can cause a computer crash. So, be careful. As a matter of fact, ISEC will throw a warning if it detects a file is binary to prevent such a matter.

### write()

- **Purpose**
    Used to create files with the contents given by the user

- **Syntax**
    
    ```
    write(filename, content)
    ```
    -  ***filename***
        This is the name that the file will be given. Similar to ``read()`` you can define a path in which the file will be made in. For example, ``write("../../huh.txt", content)``.
    - ***content***
        This is what will be written into the file once created. You can use *Text Blocks* to quickly write what info will be in the file. For example, ``write("hello.txt", [[[ Hello ]]])``. Or, just use a variable
        ```
        x = "Hello World!"
        write("hello.txt", x)
        ```
- **Notes**
    - Unlike ``read()``, ``write()`` is standalone, and doesn't require itself to be assigned to a variable. After all, what value would even be assigned to said variable?