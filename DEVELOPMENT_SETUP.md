# Development Environment Setup Guide

Welcome to the DigiSocialBlock (Nexus Protocol) project! This guide will walk you through setting up your local development environment.

## Prerequisites

- **Git:** For version control. [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- **Go:** The primary backend language. We recommend using the latest stable version. [Install Go](https://golang.org/doc/install)
- **Protocol Buffers (Protobuf) Compiler:** For compiling our `.proto` data definitions. [Install Protobuf Compiler](https://grpc.io/docs/protoc-installation/)
- **A code editor or IDE:** We recommend VS Code with the Go extension, but feel free to use your preferred editor.

## Setup Instructions

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/your-repo/digisocialblock.git
    cd digisocialblock
    ```
    *(Note: Replace with the actual repository URL once established.)*

2.  **Install Go Dependencies:**
    Once you have Go installed, you can download the necessary dependencies for the project.
    ```bash
    go mod download
    ```
    *(Note: This command will be effective once Go modules are added to the project.)*

3.  **Install Protobuf Go Plugins:**
    You'll need the Go plugins for the Protobuf compiler.
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```
    Ensure your `$GOPATH/bin` is in your system's `PATH` to use these tools.

4.  **Editor Configuration (VS Code Example):**
    - Install the official [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go).
    - Upon opening the project, the extension may prompt you to install additional Go tools. Please do so.
    - Configure the linter and formatter to adhere to project standards (details to be added).

## Next Steps

Once your environment is set up, a good place to start is by exploring the documents in the `tech_specs/` and `implementation_plans/` directories to understand the project architecture and current goals.
