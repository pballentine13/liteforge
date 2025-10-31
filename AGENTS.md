# Agent Coding Guidelines for liteforge

This document outlines the conventions and commands for coding agents operating within the liteforge repository.

## 1. Commands

| Action | Command | Example |
| :--- | :--- | :--- |
| **Build** | `go build ./...` | |
| **Test All** | `go test ./...` | |
| **Single Test** | `go test -run <TestName> <PackagePath>` | `go test -run TestNewDatastore ./internal/orm` |
| **Format** | `go fmt ./...` | |
| **Lint/Vet** | `go vet ./...` | |

## 2. Code Style & Conventions

*   **Formatting:** All code must be formatted using `gofmt`.
*   **Imports:** Imports should be grouped: standard library, then third-party, separated by a blank line.
*   **Naming:** Follow standard Go conventions: `camelCase` for unexported identifiers, `PascalCase` for exported identifiers. Acronyms (e.g., ID, HTTP) should be all caps.
*   **Error Handling:** Errors must be the last return value. Check errors immediately. Use `fmt.Errorf` for error wrapping.
*   **Documentation:** All exported functions, types, and variables must have clear, concise doc comments.
*   **Structure:** Keep functions short and focused. Prefer small interfaces.
