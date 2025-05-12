# BugLang Documentation

1. Introduction

BugLang is a simple, object-oriented programming language designed for ease of use and flexibility. It supports classes, inheritance, arrays, and basic control structures, making it suitable for both beginners and advanced developers. BugLang includes built-in libraries like `math` for common operations and provides a straightforward syntax for rapid development.

## Installation

*Note: Installation instructions are not yet available as BugLang is under development. To run BugLang code, you need a compatible interpreter or compiler (to be provided). Check back for updates on setup instructions.*

## Syntax and Features

### 1. Classes and Inheritance

BugLang supports object-oriented programming with classes and multiple inheritance. Classes are defined using the `class` keyword, and inheritance is specified by listing parent classes in parentheses.

- **Class Definition**:

  ```buglang
  class A() {
      func init() {
          this.name = "A";
      }
      func getName() {
          return this.name;
      }
  }
  ```

- **Inheritance**: Use `super(ClassName)` to call methods from a specific parent class.

  ```buglang
  class C(B, A) {
      func init() {
          super(A).init();
          println(this.name); // Outputs: A
          super(B).init();
          println(this.name); // Outputs: B
          this.name = "C";
      }
  }
  ```

### 2. Arrays and Lists

BugLang supports dynamic arrays (or lists) for storing collections of data. Arrays can hold mixed types, including numbers, strings, and objects.

- **Array Declaration**:

  ```buglang
  numbers = [20.4324, 324.423432, 10.21];
  ```

- **Array Methods**:

  - `add(value)`: Appends a value to the array.
  - `size()`: Returns the length of the array.

  ```buglang
  numbers.add(20.1122);
  ```

### 3. Loops

BugLang provides a `for` loop for iteration. The loop condition is specified without parentheses, and the loop body is executed while the condition is true.

- **For Loop**:

  ```buglang
  for (i < users.size()) {
      user = users[i];
      println("name: ", user["name"]);
      i = i + 1;
  }
  ```

### 4. Built-in Functions

- `println(...)`: Outputs arguments to the console, with optional concatenation.

- `math.round(number, decimals)`: Rounds a number to the specified decimal places.

  ```buglang
  println(math.round(20.4324, 2)); // Outputs: 20.43
  ```

### 5. Objects and Dictionaries

BugLang supports dictionary-like objects for key-value pairs, accessible using square bracket notation (`user["key"]`).

- **Example**:

  ```buglang
  users = [
      { "name": "Samandar", "age": 20, "id": 1 },
      { "name": "Nomalum", "age": 100, "id": 2 }
  ];
  ```

## Example Code

Below is a complete example demonstrating BugLang's features:

```buglang
import "math";

// Define class A
class A() {
    func init() {
        this.name = "A";
    }
    func getName() {
        return this.name;
    }
}

// Define class B
class B() {
    func init() {
        this.name = "B";
    }
    func getName() {
        return this.name;
    }
}

// Define class C inheriting from B and A
class C(B, A) {
    func init() {
        super(A).init();
        println(this.name); // Outputs: A
        super(B).init();
        println(this.name); // Outputs: B
        this.name = "C";
    }
    func getName() {
        return this.name;
    }
}

// Create instance of C
c = new C();
println(c.getName()); // Outputs: C

// Array of user objects
users = [
    { "name": "Samandar", "age": 20, "id": 1 },
    { "name": "Nomalum", "age": 100, "id": 2 }
];

// Iterate over users
i = 0;
for (i < users.size()) {
    user = users[i];
    println();
    println("id: ", user["id"]);
    println("name: ", user["name"]);
    println("age: ", user["age"]);
    println();
    i = i + 1;
}

// Array of numbers
numbers = [20.4324, 324.423432, 10.21];
numbers.add(20.1122);

// Round and print numbers
i = 0;
for (i < numbers.size()) {
    println(math.round(numbers[i], 2));
    i = i + 1;
}
```

### Output

```
A
B
C

id: 1
name: Samandar
age: 20

id: 2
name: Nomalum
age: 100

20.43
324.42
10.21
20.11
```

## Contributing

Contributions to BugLang are welcome! Please submit issues or pull requests to the repository (link TBD). Ensure your code follows the language's style and includes tests.

## License

BugLang is licensed under the MIT License. See the LICENSE file for details (to be added).