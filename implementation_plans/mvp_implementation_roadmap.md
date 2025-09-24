# **DigiSocialBlock: MVP Implementation Roadmap**

## **Objective**

The objective of this roadmap is to outline the engineering tasks required to build a Minimum Viable Product (MVP) for the DigiSocialBlock's backend infrastructure. The MVP will consist of core Go services for identity and content management, a stubbed P2P communication layer, and a unified API gateway for client interaction.

---

## **Phase 1: Foundational Setup & Core Identity Service**

*   **Goal:** To establish a working development environment and build the core service for creating, managing, and authenticating sovereign user identities.

*   **Sprint 1.1: Environment & Tooling**
    *   **Task 1.1.1:** Finalize the `DEVELOPMENT_SETUP.md` document with detailed instructions for setting up Go, Protobuf compilers (protoc), and required Go libraries.
    *   **Task 1.1.2:** Implement a Makefile or build script to automate common tasks like compiling protobufs, running tests, and building binaries.
    *   **Task 1.1.3:** Configure a linter (e.g., `golangci-lint`) and code formatter (`gofmt`) to enforce code quality standards.

*   **Sprint 1.2: Protobuf Compilation**
    *   **Task 1.2.1:** Generate the Go code from our `.proto` files (`nexus_objects.v1.proto`). The generated code will be placed in the `pkg/proto/` directory.
    *   **Task 1.2.2:** Create a simple test to ensure the generated Go structs can be created and serialized correctly.

*   **Sprint 1.3: Identity Service - Core Logic**
    *   **Task 1.3.1:** Implement the `pkg/identity` service. This service will manage the lifecycle of `NexusUserObjectV1`.
    *   **Task 1.3.2:** Implement the user creation logic: accept a public key and handle, generate a `user_id`, and store the new user object (initially in-memory or a simple file-based DB).
    *   **Task 1.3.3:** Implement a method to retrieve a user object by `user_id` or `handle`.

*   **Sprint 1.4: Identity Service - Authentication**
    *   **Task 1.4.1:** Implement a cryptographic utility package for signing and verifying messages using the Ed25519 standard.
    *   **Task 1.4.2:** Implement an authentication challenge-response mechanism. A client must sign a server-provided nonce to prove ownership of their private key.

---

## **Phase 2: Content Service & Initial EchoNet Stub**

*   **Goal:** To build the service responsible for content management and to create a stubbed implementation of the EchoNet for initial integration testing.

*   **Sprint 2.1: Content Service - Core Logic**
    *   **Task 2.1.1:** Implement the `pkg/content` service to manage `NexusContentObjectV1`.
    *   **Task 2.1.2:** Implement the content creation logic: accept a request containing an `author_user_id`, `content_body_uri`, etc., validate the author's signature, and store the content object.
    *   **Task 2.1.3:** Implement methods to retrieve content objects by `content_id` and to get all content by a specific `author_user_id`.

*   **Sprint 2.2: EchoNet P2P Stub**
    *   **Task 2.2.1:** Implement the `pkg/echonet/core` service as a stubbed interface.
    *   **Task 2.2.2:** Create a `Publish` method that, for the MVP, simply logs the content that *would* be broadcast to the P2P network.
    *   **Task 2.2.3:** Integrate this `Publish` call into the content service upon successful creation of a new content object. This sets the stage for future P2P implementation.

---

## **Phase 3: Backend Integration & API Gateway**

*   **Goal:** To tie the identity and content services together and expose their functionality through a single, unified API gateway.

*   **Sprint 3.1: Service Integration**
    *   **Task 3.1.1:** Ensure the content service properly calls the identity service to authenticate users before processing requests.
    *   **Task 3.1.2:** Refine data storage to use a more robust embedded database (e.g., BoltDB or Badger) instead of in-memory maps.

*   **Sprint 3.2: gRPC API Gateway**
    *   **Task 3.2.1:** Define a new `.proto` file for the API gateway service, specifying the RPC methods for user creation, login, content posting, etc.
    *   **Task 3.2.2:** Implement a gRPC server that acts as the API gateway. This server will orchestrate calls to the underlying identity and content services.
    *   **Task 3.2.3:** Implement end-to-end tests for the primary user flows through the gRPC gateway.

*   **MVP Completion:** At the end of Phase 3, we will have a testable backend MVP with core services for identity and content, and a gRPC API for client applications to interact with.