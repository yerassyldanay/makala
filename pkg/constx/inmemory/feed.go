package inmemory

import "fmt"

// FeedVersion
// simplifies the way of creating key for in-memory databases
type FeedVersion string

var (
	FeedVersionOld   FeedVersion = "old"
	FeedVersionNew   FeedVersion = "new"
	FeedVersionTemp  FeedVersion = "temporary"
	FeedVersionEmpty FeedVersion = "empty"
)

// GetKey returns a key
func (f FeedVersion) GetKey() string {
	return fmt.Sprintf("key_feed_version_%s", f)
}
