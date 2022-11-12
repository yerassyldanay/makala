package postfeed

import (
	"context"

	"github.com/yerassyldanay/makala/provider/poststore"
	"github.com/yerassyldanay/makala/server/rest/model"
	model2 "github.com/yerassyldanay/makala/service/postfeed/model"
)

type Poster interface {
	// CreatePost stores a post on the datastore
	CreatePost(ctx context.Context, args poststore.CreateParams) (poststore.FeedPost, error)
	// GetChunkByChunk
	// fetches posts chunk by chunk without any ranking by score
	GetChunkByChunk(ctx context.Context, page, count int32) ([]poststore.FeedPost, error)
	// GetFeed
	// returns n number of best-ranked posts (makala feed), including promoted ones
	GetFeed(ctx context.Context, args model2.GetFeedParams) (model.GetFeedResponse, error)
}
