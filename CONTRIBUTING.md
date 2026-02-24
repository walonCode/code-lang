# Contributing to Code-Lang 🚀

First off, thank you for considering contributing to Code-Lang! It's people like you that make Code-Lang a better tool for everyone.

This project is a passion project built in Go, and we welcome all types of contributions—from bug reports and documentation improvements to new standard library modules and core language features.

---

## 🛠️ Development Setup

### Prerequisites
- **Go**: Version 1.22 or higher is recommended.
- **Git**: For version control.

   3. Verify you can build the project:
   
   ./code-lang --version
   

## 🏗️ Project Structure

Understanding the layout of the project will help you find where to make changes:

- `lexer/`: Lexical analysis (converts source code to tokens).
- `token/`: Token definitions.
- `parser/`: Recursive descent parser (converts tokens to AST).
- `ast/`: Abstract Syntax Tree nodes.
- `object/`: The type system and object represention used during evaluation.
- `evaluator/`: The core logic that traverses the AST and executes it.
- `std/`: The standard library modules (written in Go).

---

## 🧪 Testing

We value stability. If you add a new feature or fix a bug, please include tests.

Run all tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./evaluator
```

---

## 📚 Adding to the Standard Library

Adding a new module (like `math` or `strings`) is a great way to contribute.

1. Create a new directory and file in `std/`: `std/your_module/your_module.go`.
2. Implement a `Module()` function that returns an `*object.Module`.
3. Register your module in `evaluator/builtin.go` inside the `init()` function.

**Example Module Structure:**
```go
package your_module

func Module() *object.Module {
    return &object.Module{
        Members: map[string]object.Object{
            "hello": &object.Builtin{
                Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
                    return &object.String{Value: "world"}
                },
            },
        },
    }
}
```

---

## 📝 Pull Request Process

1. Create a new branch for your feature or fix: `git checkout -b feat/my-new-feature`.
2. Commit your changes with clear, descriptive messages.
3. Push to your fork and submit a Pull Request.
4. Provide a clear description of the change and link any related issues.

---

## 💬 Communication

If you have questions or want to discuss a major change before writing code, feel free to open a **GitHub Discussion** or an **Issue**.

Happy coding!
