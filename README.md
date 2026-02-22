# Code-Lang 🚀

Code-Lang is a modern, interpreted programming language written in Go. It began as an implementation following the excellent book **"Writing An Interpreter In Go"** by [Thorsten Ball](https://interpreterbook.com/), and has since evolved with additional features and custom extensions.

> [!IMPORTANT]
> **Status:** Code-Lang is a **passion project** and is currently under active development. While the core language is functional, it is not production-ready. If you intend to use this in a production environment, significant work, security audits, and optimizations are required.

---

## ✨ Features

- **Rich Type System:**
  - Integers and Floats
  - Strings and Characters
  - Booleans
  - Arrays (e.g., `[1, 2, 3]`)
  - Hashes/Dictionaries (e.g., `{"name": "Code-Lang"}`)
- **First-Class Functions:** Function literals, closures, and higher-order functions.
- **Control Flow:**
  - `if-elseif-else` expressions (everything is an expression!).
  - `while` loops for simple iteration.
  - `for` loops for structured iteration.
- **Support for Comments:** Single-line (`#`) and multi-line (`/* */`).
- **Standard Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`, `%` (Modulo)
  - Advanced: `**` (Power), `//` (Floor Division), `=` (Assignment)
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `!` (Negation)
- **Built-in Functions:** `print`, `printf`, `typeof`, `len`, `push`, and more.
- **Module System:** Import other `.cl` files using `import "module"`.
- **Member Access:** Dot notation (`obj.prop`) for Hashes and Modules.
- **Compound Assignment:** Supports `+=`, `-=`, `*=`, `/=`, etc.
- **REPL:** Interactive shell with precise line/column error tracking.
- **File Execution:** Run scripts with the `.cl` extension.

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.25.3 or higher recommended)

### Installation

#### Option 1: Using `go install` (Recommended)
You can install the **code-lang** binary directly to your `$GOPATH/bin`:

```bash
go install github.com/walonCode/code-lang@latest
```

#### Option 2: Pre-built Binaries
Head over to the [Releases](https://github.com/walonCode/code-lang/releases) section to download path-ready binaries for Windows, macOS, and Linux.

#### Option 3: From Source
Clone and build manually:

```bash
git clone https://github.com/walonCode/code-lang.git
cd code-lang
go build -o code-lang main.go
```

### Running the REPL

Start the interactive shell by running:

```bash
go run main.go
```

### Running a Script

You can execute a Code-Lang script by passing the filename as an argument:

```bash
go run main.go hello.cl
```

---

## 📖 Language Syntax at a Glance

### Variables & Functions

```rust
let age = 25;
let name = "Developer";
let isLearning = true;

let add = fn(a, b) {
    return a + b;
};

print(add(10, 15)); // Output: 25
```

### Arrays and Indexing

```rust
let fibonacci = [0, 1, 1, 2, 3, 5, 8];


let person = {"name": "Alice", "age": 30};

```

### Conditionals

```rust
let x = 10;
let result = if (x > 10) {
    "Greater"
} elseif (x == 10) {
    "Equal"
} else {
    "Smaller"
};
```

### Loops

```rust
// While loop
let i = 0;
while (i < 5) {
    print(i);
    i = i + 1;
};

// For loop
for (let j = 0; j < 5; j += 1) {
    print(j);
};
```

### Modules & Member Access

```rust
// math_lib.cl
let PI = 3.14159;
let square = fn(x) { x * x; };

// main.cl
import "math_lib";
print(math_lib.PI);
print(math_lib.square(10));

// Hashes
let user = {"name": "Thorsten", "active": true};
user.name = "Walon";
user.score = 100;
user.score += 50;
print(user.name); // Walon
print(user.score); // 150
```

### Advanced Features

```rust
// Comments
# This is a single line comment
/* 
   This is a 
   multi-line comment 
*/

// Formatted print
let name = "Alice";
printf("Hello, %s!", name);

// Type checking
print(typeof(10)); // Output: INTEGER
print(typeof("hi")); // Output: STRING
```

---

## 🗺 Roadmap

We are constantly working to make Code-Lang better. Here is what's coming next:

- [x] **Better Error Reporting:** Line and column tracking for precise debugging.
- [x] **Comments:** Support for single and multi-line comments.
- [x] **Loops:** Implementing `while` and `for` loops.
- [ ] **Logical Operators:** Adding `&&` (AND) and `||` (OR) with short-circuiting.
- [ ] **Standard Library (Internal):** Dedicated Go-backed modules for `math`, `fs`, and `http`.
- [x] **Import System:** Ability to include other `.cl` files.
- [x] **Member Access:** Dot notation for objects and modules.
- [x] **Compound Assignment:** Support for `+=`, `-=`, etc.

---

## 📜 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Thorsten Ball** for the foundational guide [Writing An Interpreter In Go](https://interpreterbook.com/).
- The Go community for providing an incredible ecosystem for language development.
