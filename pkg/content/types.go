package content

import (
	"time"
)

// ContentObjectV1 represents a single, immutable piece of content.
// This is the native Go data structure for content.
type ContentObjectV1 struct {
	// A unique, immutable identifier for the content object (e.g., a UUID).
	ContentID string

	// The user_id of the author. This links the content to its sovereign creator.
	AuthorUserID string

	// The timestamp when the content was created.
	CreatedAt time.Time

	// A URI pointing to the actual content body, which is stored off-chain.
	ContentBodyURI string

	// The SHA-256 hash of the off-chain content body, ensuring data integrity.
	ContentBodyHash []byte

	// The author's signature over the ContentBodyHash, proving authorship.
	Signature []byte
}