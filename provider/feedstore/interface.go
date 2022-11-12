package feedstore

import (
	"context"
	"time"

	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
)

// Storer abstracts all methods for working with feed
type Storer interface {
	// CreatePost accepts the version of a feed & post & adds it to sorted set,
	// which will be ranked by their scores
	CreatePost(ctx context.Context, version inmemory.FeedVersion, score float64, id int64) error
	// RemovePostsWithMinScore leaves 'maxElements' posts in a feed,
	// other posts will be removed
	RemovePostsWithMinScore(ctx context.Context, feedVersion inmemory.FeedVersion, maxElements int64) error
	// CopyFeed copies feed from one version to another,
	// while the initial posts in destination feed will be removed
	CopyFeed(ctx context.Context, fromThis, toThis inmemory.FeedVersion) error
	// GetFeed returns 'count' number of posts starting from 'page * count'
	GetFeed(ctx context.Context, feedVersion inmemory.FeedVersion, page, count int64) ([]int64, error)
	// GetFeedLength returns number of posts in a given feed version
	GetFeedLength(ctx context.Context, feedVersion inmemory.FeedVersion) (int64, error)
	// SetUpdateAt sets a time when feed was recreated/updated
	SetUpdateAt(ctx context.Context, updatedAt time.Time) error
	// GetUpdateAt returns the time when the feed was updated by background job
	GetUpdateAt(ctx context.Context) time.Time
}
