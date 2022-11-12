package synchron

import (
	"context"
	"time"

	"github.com/yerassyldanay/makala/pkg/constx"
	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
	"github.com/yerassyldanay/makala/provider/poststore"
)

func (r SyncRun) UpdateFeed() {
	// to empty the temporary sorted set
	if err := r.FeedStorage.CopyFeed(context.Background(), inmemory.FeedVersionEmpty, inmemory.FeedVersionTemp); err != nil {
		r.Log.Printf("[CRON-UPDATE] failed to empty temporary sorted set. err: %v \n", err)
		return
	}

	var limit int32 = 1000
	var offset, failed int32
	for {
		postFetched, err := r.PostRepo.GetAll(context.Background(), poststore.GetAllParams{
			Offset: offset,
			Limit:  limit,
		})
		offset = offset + limit
		if err != nil && failed >= 5 {
			r.Log.Printf("[CRON-UPDATE] failed to fetch posts %d times. err: %v \n", failed, err)
			continue
		}
		if err != nil {
			failed += 1
			r.Log.Printf("[CRON-UPDATE] failed to fetch posts. err: %v \n", err)
			continue
		}

		if len(postFetched) == 0 {
			break
		}

		for _, eachPost := range postFetched {
			err = r.FeedStorage.CreatePost(context.Background(), inmemory.FeedVersionTemp, 0, eachPost.ID)
			if err != nil {
				r.Log.Printf("[CRON-UPDATE] %v \n", err)
				continue
			}
			err = r.FeedStorage.RemovePostsWithMinScore(context.Background(), inmemory.FeedVersionTemp, constx.FeedMaxLength)
			if err != nil {
				r.Log.Printf("[CRON-UPDATE] %v \n", err)
				continue
			}
		}
	}

	// copy 'new' sorted set to 'old'
	// copy newly created 'temp' sorted set to 'new'
	if err := r.FeedStorage.CopyFeed(context.Background(), inmemory.FeedVersionNew, inmemory.FeedVersionOld); err != nil {
		r.Log.Printf("[CRON-UPDATE] copy sets new -> old. err: %v \n", err)
	}
	if err := r.FeedStorage.CopyFeed(context.Background(), inmemory.FeedVersionTemp, inmemory.FeedVersionNew); err != nil {
		r.Log.Printf("[CRON-UPDATE] copy sets temp -> new. err: %v \n", err)
	}
	if err := r.FeedStorage.CopyFeed(context.Background(), inmemory.FeedVersionEmpty, inmemory.FeedVersionTemp); err != nil {
		r.Log.Printf("[CRON-UPDATE] copy sets empty -> temporary. err: %v \n", err)
	}
	if err := r.FeedStorage.SetUpdateAt(context.Background(), time.Now()); err != nil {
		r.Log.Printf("[CRON-UPDATE] failed to set update time (the last time when feed was updated). err: %v \n", err)
	}
}
