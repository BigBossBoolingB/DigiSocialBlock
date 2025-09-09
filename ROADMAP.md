# DigiSocialBlock (Nexus Protocol) - Development Roadmap

This document outlines a high-level roadmap for the phased development of the DigiSocialBlock. The milestones are designed to build upon each other, starting with core functionalities and progressively layering in the more advanced, conceptual features of the protocol.

---

## Milestone 1: Foundation & Sovereign Identity

*   **Focus:** Establish the core technical foundation and allow users to create and control their digital identity. This milestone is about setting up the board.
*   **Key Deliverables:**
    1.  **Project Setup:** Initialize a monorepo containing the client-side application (e.g., Next.js) and the smart contract development environment (e.g., Hardhat).
    2.  **Wallet Integration:** Client application can connect to user wallets (e.g., via MetaMask).
    3.  **Lens Protocol Integration:** Integrate the Lens Protocol SDKs for interacting with the Polygon testnet.
    4.  **Feature: Create Profile:** Users can create their own `.lens` Profile NFT, establishing their sovereign identity on the network.
    5.  **Feature: Manage Profile:** Users can update their profile metadata (display name, bio, avatar). Metadata will be stored on Arweave and linked from their Profile NFT.
*   **Outcome:** A basic, functional client where users can manage their on-chain identity.

---

## Milestone 2: Core Content & The Resonant Conduit (v1)

*   **Focus:** Enable the creation and consumption of content, forming the initial version of the "Resonant Conduit."
*   **Key Deliverables:**
    1.  **Arweave Integration:** Implement client-side functionality to upload content to Arweave.
    2.  **Feature: Create Post:** Users can publish a simple text-based post. The client uploads the content to Arweave and broadcasts the link via a Lens Protocol publication.
    3.  **Feature: Basic Resonance Stream:** The client can fetch and display a chronological feed of posts from profiles the user follows.
    4.  **Content Hydration:** The client correctly resolves `contentURI`s from the blockchain and loads the corresponding content from Arweave.
*   **Outcome:** Users can create and view posts, forming a simple, decentralized social feed.

---

## Milestone 3: The Immutable Social Contract (ISC)

*   **Focus:** Implement the custom, tiered social contract system that governs user interactions. This is a core innovation of the protocol.
*   **Key Deliverables:**
    1.  **ISC Smart Contracts:** Develop, test, and deploy the `ISCFactory` and `ISC` smart contracts as defined in the technical specification.
    2.  **Custom Reference Module:** Develop a custom Lens `referenceModule` smart contract that checks for an active ISC between two users before allowing an interaction (e.g., a comment).
    3.  **Feature: Manage ISCs:** Users can initiate, accept, and terminate Tier 1 ("Public Resonance") social contracts with other users.
    4.  **Feature: Gated Interactions:** The ability to comment on a post is gated by the ISC system. Only users with an active Tier 1 (or higher) contract can comment.
*   **Outcome:** The social graph moves beyond simple "follows" to explicit, auditable, on-chain relationships.

---

## Milestone 4: The Axiomatic Filter & Advanced Resonance

*   **Focus:** Enhance the user experience with client-side intelligence and more sophisticated social interactions.
*   **Key Deliverables:**
    1.  **Axiomatic Filter (v1):** Implement a client-side filter using TensorFlow.js. The initial version will use a pre-trained model to filter for common "noise" like toxicity or spam. This happens entirely on the user's device.
    2.  **Frequencies of Resonance:** The client will allow users to filter their Resonance Stream by the tags (`#Strategy`, `#Poetry`, etc.) defined in the content's Arweave body.
    3.  **Advanced ISC Tiers:** Implement the logic for Tier 2 ("Collaborative Synergy") and Tier 3 ("The Inner Circle"), which could unlock features like private messaging or access to shared spaces (to be defined in a future spec).
*   **Outcome:** A more intelligent and customizable user experience that actively protects the user's "conceptual environment."

---

## Milestone 5: Anti-Fragile Governance

*   **Focus:** Begin the process of decentralizing the protocol's governance, fulfilling the vision of an anti-fragile, community-driven system.
*   **Key Deliverables:**
    1.  **Governance Tokenomics:** A proposal for the protocol's governance token, including distribution and utility.
    2.  **DAO Smart Contracts:** Deployment of a basic DAO framework (e.g., based on Governor Bravo) that allows token holders to vote on proposals.
    3.  **Governance UI:** A section in the client application where users can view, create, and vote on governance proposals related to the protocol's evolution.
*   **Outcome:** The protocol has a clear path toward community ownership and decentralized decision-making.
