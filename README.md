# Nbx

A collection of reusable packages and CLI tools for common development tasks like logging, database utilities, data conversion, transaction management, and enum generation.

## Introduction

`nbx` is a personal Go toolbox containing a collection of reusable packages and command-line interface (CLI) tools. The goal is to provide convenient, well-tested components for common tasks encountered during Go development, fostering code reuse and efficiency.

## Features

The toolbox currently includes the following key packages and CLI tools:

* **`pkg/logger`**: A flexible logging library with different modes, adapters (e.g., gRPC interceptors), encoders, and configurable log levels.
* **`pkg/dbu`**: Database utilities including helpers for building update maps for relational databases (RDB, like GORM) and MongoDB, handling struct tags, default naming, and pointer values.
* **`pkg/pbconv`**: Protocol Buffer conversion helpers (details based on code, assumed to help convert between Go types and protobuf types).
* **`pkg/transaction`**: Helpers and a `TransactionManager` pattern for managing database transactions cleanly, especially with GORM, by passing the transaction via `context.Context`.
* **`cmd/enumgo`**: A CLI tool (`nbx enumgo`) to generate Go enum code directly from YAML definition files.

## Getting Started

To get the `nbx` toolbox and its source code, you can clone the repository:

```bash
git clone [https://github.com/byte4cat/nbx.git](https://github.com/yimincai/nbx.git)
```
