package synchron

import (
	"log"

	"github.com/yerassyldanay/makala/provider/feedstore"
	"github.com/yerassyldanay/makala/provider/poststore"
)

type SyncRun struct {
	Log         *log.Logger
	PostRepo    poststore.Querier
	FeedStorage feedstore.Storer
}

func NewSyncRunner(log *log.Logger, postRepo poststore.Querier, feedStorage feedstore.Storer) *SyncRun {
	return &SyncRun{
		Log:         log,
		PostRepo:    postRepo,
		FeedStorage: feedStorage,
	}
}
