# **DigiSocialBlock: Development Environment Setup**

## **1. Overview**

This document provides the canonical instructions for setting up a local development environment to contribute to the DigiSocialBlock project. Adhering to these steps ensures a consistent and stable environment, aligning with our **Expanded KISS Principle** of clarity and synchronization.

---

## **2. Prerequisites**

You will need the following tools installed on your system:
*   **Go** (version 1.21 or later)
*   **Protocol Buffer Compiler** (`protoc` version 3.x)
*   **Git** for version control
*   A text editor or IDE, **Visual Studio Code** is recommended.

---

## **3. Go Installation**

1.  **Download and Install Go:** Visit the official Go website at [go.dev/dl/](https://go.dev/dl/) and download the installer for your operating system (macOS, Linux, Windows).
2.  **Follow Installation Instructions:** Run the installer and follow the on-screen prompts.
3.  **Verify Installation:** Open a new terminal or command prompt and run the following command. You should see the installed Go version.
    ```sh
    go version
    ```
4.  **Configure GOPATH:** Ensure your `GOPATH` environment variable is set up correctly. This is where Go workspaces and downloaded packages reside. For most users, this will be set up automatically. You can check its location with `go env GOPATH`.

---

## **4. Protocol Buffers Setup**

The Protocol Buffer compiler (`protoc`) is required to generate Go code from our `.proto` definitions.

### **4.1. Install `protoc` Compiler**

*   **macOS (using Homebrew):**
    ```sh
    brew install protobuf
    ```
*   **Linux (using `apt`):**
    ```sh
    sudo apt-get install -y protobuf-compiler
    ```
*   **Windows (using `winget` or manually):**
    ```sh
    winget install Google.Protobuf
    ```
    Alternatively, download the pre-compiled `protoc-*.zip` from the [Protobuf GitHub Releases](https://github.com/protocolbuffers/protobuf/releases) page, extract it, and add the `bin` directory to your system's `PATH`.

*   **Verify Installation:**
    ```sh
    protoc --version
    ```
    This should output `libprotoc` followed by the version number (e.g., `libprotoc 3.20.1`).

### **4.2. Install Go Plugins for `protoc`**

These plugins are used by `protoc` to generate Go-specific code.

1.  **`protoc-gen-go` (for standard Protobuf messages):**
    ```sh
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    ```
2.  **`protoc-gen-go-grpc` (for gRPC services):**
    ```sh
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    ```
3.  **Ensure Plugins are in PATH:** The `go install` command will place the binaries in your `$GOPATH/bin` directory. Make sure this directory is included in your system's `PATH` environment variable.

---

## **5. Project Setup**

1.  **Clone the Repository:**
    ```sh
    git clone [repository-url]
    cd [repository-directory]
    ```
2.  **Install Go Dependencies:** Navigate to the `pkg/` directory and install the necessary Go modules defined in `go.mod`.
    ```sh
    cd pkg
    go mod tidy
    cd ..
    ```
    The `go mod tidy` command ensures you have all the required dependencies for the backend services.

---

## **6. Recommended IDE Setup (Visual Studio Code)**

1.  **Install VS Code:** Download from [code.visualstudio.com](https://code.visualstudio.com/).
2.  **Install the Go Extension:** Open VS Code, go to the Extensions view (Ctrl+Shift+X), and search for and install the official `Go` extension by the Go Team at Google.
3.  **Configure Workspace Settings:** Create a `.vscode` directory in the project root and add a `settings.json` file with the following content. This will enable formatting on save and ensure consistent linting.

    **`.vscode/settings.json`:**
    ```json
    {
      "go.useLanguageServer": true,
      "editor.formatOnSave": true,
      "[go]": {
        "editor.defaultFormatter": "golang.go"
      },
      "gopls": {
        "ui.semanticTokens": true
      }
    }
    ```
    VS Code may prompt you to install additional Go tools (like `gopls` and `golangci-lint`) when you first open a `.go` file. Click "Install All" to proceed.

---

## **7. Next Steps**

With the environment set up, you are now ready to contribute. Consult the `implementation_plans/mvp_implementation_roadmap.md` for the current development phase and tasks. To compile the `.proto` files, you can use the `make proto` command (once the `Makefile` is implemented).