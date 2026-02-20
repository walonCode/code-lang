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
- **Control Flow:** `if-else` expressions (everything is an expression!).
- **Standard Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Advanced: `**` (Power), `//` (Floor Division)
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `!` (Negation)
- **Built-in Functions:** `print`, `len`, `first`, `last`, `rest`, `push`, and more.
- **REPL:** Interactive shell with a friendly greet.
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
let result = if (x > 5) {
    "Greater"
} else {
    "Smaller or Equal"
};
```

---

## 🗺 Roadmap

We are constantly working to make Code-Lang better. Here is what's coming next:

- [ ] **Better Error Reporting:** Line and column tracking for precise debugging.
- [ ] **Comments:** Support for `//` and `/* */`.
- [ ] **Loops:** Implementing `while` and `for` loops.
- [ ] **Logical Operators:** Adding `&&` (AND) and `||` (OR) with short-circuiting.
- [ ] **Standard Library:** Dedicated modules for `math`, `fs`, and `http`.
- [ ] **Import System:** Ability to include other `.cl` files.

---

## 📜 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Thorsten Ball** for the foundational guide [Writing An Interpreter In Go](https://interpreterbook.com/).
- The Go community for providing an incredible ecosystem for language development.
