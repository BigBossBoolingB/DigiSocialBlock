# **DigiSocialBlock: Development Environment Setup**

## **1. Overview**

This document provides the canonical instructions for setting up a local development environment to contribute to the DigiSocialBlock project. Adhering to these steps ensures a consistent and stable environment.

---

## **2. Prerequisites**

You will need the following tools installed on your system:
*   **Go** (version 1.21 or later)
*   **Git** for version control
*   A text editor or IDE, **Visual Studio Code** is recommended.

---

## **3. Go Installation**

1.  **Download and Install Go:** Visit the official Go website at [go.dev/dl/](https://go.dev/dl/) and download the installer for your operating system.
2.  **Verify Installation:** Open a new terminal and run the following command. You should see the installed Go version.
    ```sh
    go version
    ```
3.  **Configure GOPATH:** Ensure your `GOPATH` environment variable is set up correctly. For most users, this will be set up automatically. You can check its location with `go env GOPATH`.

---

## **4. Project Setup**

1.  **Clone the Repository:**
    ```sh
    git clone [repository-url]
    cd [repository-directory]
    ```
2.  **Install Go Dependencies:** Navigate to the `pkg/` directory (once created) and install the necessary Go modules defined in `go.mod`.
    ```sh
    cd pkg
    go mod tidy
    cd ..
    ```

---

## **5. Recommended IDE Setup (Visual Studio Code)**

1.  **Install the Go Extension:** Open VS Code, go to the Extensions view (Ctrl+Shift+X), and install the official `Go` extension by the Go Team at Google.
2.  **Configure Workspace Settings:** Create a `.vscode` directory in the project root and add a `settings.json` file with the following content to enable formatting on save.

    **`.vscode/settings.json`:**
    ```json
    {
      "go.useLanguageServer": true,
      "editor.formatOnSave": true,
      "[go]": {
        "editor.defaultFormatter": "golang.go"
      }
    }
    ```
    VS Code may prompt you to install additional Go tools (like `gopls`). Click "Install All" to proceed.

---

## **6. Next Steps**

With the environment set up, you are now ready to contribute. Run `make test` to verify the project's state.