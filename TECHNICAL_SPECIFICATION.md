# DigiSocialBlock (Nexus Protocol) - Technical Specification

This document provides the technical specification for the DigiSocialBlock. It is based on the conceptual blueprint and the technology recommendations outlined in `TECHNOLOGY_CHOICES.md`.

---

## 1. System Architecture

The DigiSocialBlock is a decentralized system composed of three main layers: the Client Application, the Blockchain Layer (Polygon), and the Storage Layer (Arweave).

```
+--------------------------+
|   Client Application     |
| (Web/Desktop/Mobile)     |
| - UI/UX                  |
| - Axiomatic Filter (TF.js)|
+--------------------------+
           |
           | (Reads/Writes)
           v
+-------------------------------------------------------------+
| Blockchain Layer (Polygon)                                  |
|                                                             |
| +-----------------+  +-----------------+  +-----------------+ |
| | Lens Protocol   |  | Custom ISC      |  | Other Contracts | |
| | (Profiles,      |  | (Tiered Access  |  | (e.g., DAOs)    | |
| |  Follows, Posts)|  |  Contracts)     |  |                 | |
| +-----------------+  +-----------------+  +-----------------+ |
+-------------------------------------------------------------+
           |
           | (Content Links)
           v
+--------------------------+
|   Storage Layer (Arweave)|
| (Permanent, Immutable    |
|  Content Storage)        |
+--------------------------+
```

**Flow of Interaction:**

1.  **User Identity:** A user's identity is a **Lens Protocol Profile NFT** on the Polygon blockchain. This NFT is owned and controlled by the user's wallet.
2.  **Content Creation:** When a user posts, the client application uploads the content (text, image, etc.) to **Arweave**. Arweave returns a permanent transaction ID.
3.  **Broadcasting:** The client then creates a "Publication" via a **Lens Protocol smart contract** on Polygon. This on-chain publication does not contain the content itself, but rather the **Arweave URI** (`ar://...`) pointing to it.
4.  **Social Contracts:** User-to-user relationships are managed by custom **Immutable Social Contract (ISC) smart contracts**, also on Polygon. These contracts define the "Resonance Tier" and gate access to certain actions (e.g., commenting, viewing private content).
5.  **Content Consumption (Resonance Stream):** To view their feed, the client fetches the list of publications from profiles the user follows (via Lens Protocol). It then resolves the Arweave URIs to retrieve the actual content. Finally, the client-side **Axiomatic Filter** processes this content before it is displayed to the user.

---

## 2. Data Models

### 2.1. User Profile (Lens Protocol NFT)

A user's profile is a standard ERC-721 NFT managed by the Lens Protocol contracts.

| Field | Type | Description |
| :--- | :--- | :--- |
| `profileId` | `uint256` | The unique token ID of the profile NFT. |
| `owner` | `address` | The wallet address that owns and controls the profile. |
| `handle` | `string` | The unique, human-readable username (e.g., `architect.lens`). |
| `metadata` | `string` | A URI (typically pointing to Arweave/IPFS) for a JSON file containing off-chain profile details like display name, bio, and avatar. |

### 2.2. Content Publication (Lens Protocol)

All content posts are "Publications" created via Lens Protocol.

| Field | Type | Description |
| :--- | :--- | :--- |
| `publicationId` | `uint256` | Unique ID for the publication within a profile. |
| `profileId` | `uint256` | The ID of the profile that created the publication. |
| `contentURI` | `string` | **The Arweave URI (`ar://...`) pointing to the off-chain content.** |
| `referenceModule` | `address` | A smart contract that dictates who can comment on or mirror this post. **This will be linked to our ISCs.** |
| `collectModule` | `address` | A smart contract that dictates who can "collect" (mint an NFT of) this post. |

### 2.3. Content Body (Stored on Arweave)

The actual content of a post, stored as a JSON object on Arweave.

| Field | Type | Description |
| :--- | :--- | :--- |
| `version` | `string` | The version of the DigiSocialBlock spec (e.g., "1.0"). |
| `mainContent` | `string` | The primary content of the post (e.g., Markdown text). |
| `type` | `string` | MIME type of the main content (e.g., `text/markdown`, `image/png`). |
| `timestamp` | `uint64` | Unix timestamp of creation. |
| `tags` | `string[]` | Array of tags or "Frequencies of Resonance" (e.g., `["#Strategy", "#AI"]`). |
| `metadata` | `object` | Any other arbitrary metadata. |

### 2.4. Immutable Social Contract (ISC)

A custom smart contract to govern tiered relationships. A factory contract will be used to deploy new ISCs.

| Field | Type | Description |
| :--- | :--- | :--- |
| `partyA` | `address` | The wallet address of the user initiating the contract. |
| `partyB` | `address` | The wallet address of the user receiving the contract offer. |
| `tier` | `uint8` | **1:** Public Resonance, **2:** Collaborative Synergy, **3:** Inner Circle. |
| `status` | `uint8` | **0:** Pending, **1:** Active, **2:** Terminated. |
| `createdAt` | `uint256` | Block timestamp of contract creation. |
| `activatedAt` | `uint256` | Block timestamp of contract activation. |

---

## 3. Core API & Smart Contract Interactions

This section describes the primary functions the client will invoke. These are not direct API calls but interactions with the Polygon smart contracts.

### 3.1. User & Profile Management

*   `createProfile(handle)`: Calls the Lens Protocol hub contract to mint a new Profile NFT for the user.
*   `setProfileMetadata(profileId, metadataURI)`: Updates the profile's metadata URI to point to a new JSON file on Arweave.

### 3.2. Content & Publication

*   `createPost(profileId, contentBody)`:
    1.  Client serializes `contentBody` JSON and uploads it to Arweave, receiving an `arweave_tx_id`.
    2.  Client constructs the `contentURI` as `ar://{arweave_tx_id}`.
    3.  Client calls the Lens Protocol `post` function with the `profileId` and `contentURI`.

### 3.3. Immutable Social Contract (ISC) Management

*   `initiateISC(targetUserAddress, tier)`: The client calls our `ISCFactory` contract. The factory deploys a new `ISC` contract with `partyA` as the caller, `partyB` as the target, the specified `tier`, and `status` as `Pending`.
*   `acceptISC(iscContractAddress)`: The user who is `partyB` in the contract calls this function on the `ISC` contract. It changes the `status` to `Active` and sets the `activatedAt` timestamp.
*   `terminateISC(iscContractAddress)`: Either party can call this to change the `status` to `Terminated`.
*   `getISCStatus(userA, userB)`: A read-only function to check if an active ISC exists between two users and return its tier. This will be used by other contracts (like the `referenceModule`) to gate access.

### 3.4. Resonance Stream (Data Fetching)

*   `getResonanceStream()`: This is a client-side orchestration process:
    1.  Fetch all profiles the user follows using the Lens API.
    2.  Fetch all active ISCs for the user from our contracts.
    3.  Use the Lens API to query for recent publications from this combined list of profiles.
    4.  For each publication, fetch the content from Arweave using its `contentURI`.
    5.  Pass the fetched content through the local **Axiomatic Filter**.
    6.  Render the filtered, hydrated content to the user.
