package postfeed

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/yerassyldanay/makala/pkg/configx"
	"github.com/yerassyldanay/makala/pkg/constx"
	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
	"github.com/yerassyldanay/makala/provider/adstore"
	"github.com/yerassyldanay/makala/provider/feedstore"
	"github.com/yerassyldanay/makala/provider/poststore"
	"github.com/yerassyldanay/makala/server/rest/model"
	model2 "github.com/yerassyldanay/makala/service/postfeed/model"
)

type PostService struct {
	Conf        configx.Configuration
	Log         *log.Logger
	PostStore   poststore.Querier
	FeedStorage feedstore.Storer
	AdStorage   adstore.Storer
}

var _ Poster = PostService{}

func NewPostService(log *log.Logger, post poststore.Querier, feedStorage feedstore.Storer, adStorage adstore.Storer) *PostService {
	return &PostService{
		Log:         log,
		PostStore:   post,
		FeedStorage: feedStorage,
		AdStorage:   adStorage,
	}
}

func (s PostService) CreatePost(ctx context.Context, args poststore.CreateParams) (poststore.FeedPost, error) {
	postCreated, err := s.PostStore.Create(ctx, args)
	if err != nil {
		return poststore.FeedPost{}, err
	}

	// this is needed to consider edge cases: when there are not enough posts in feed
	// Note: (in case feed at its max) we cannot just add to the cache
	// in order to avoid the case when use sees a post more than once
	count, err := s.FeedStorage.GetFeedLength(ctx, inmemory.FeedVersionNew)
	if err != nil {
		s.Log.Printf("[FEED-STORAGE] failed to count the number of ads in cache. err: %v \n", err)
	}

	// feed - contains only not promoted posts
	// ads - only promoted posts
	if !postCreated.Promoted && count < constx.FeedMaxLength {
		err = s.FeedStorage.CreatePost(ctx, inmemory.FeedVersionNew, postCreated.Score, postCreated.ID)
		if err != nil {
			s.Log.Printf("[FEED-STORAGE] failed to add post to storage. err: %v \n", err)
		}
	} else if postCreated.Promoted {
		// if the post is promoted then add it to the ads
		if err := s.AdStorage.CreateAd(ctx, postCreated.ID); err != nil {
			s.Log.Printf("[ADS-STORAGE] failed to add a promoted post to the ads list. err: %v \n", err)
		}
	}

	return postCreated, nil
}

func (s PostService) GetChunkByChunk(ctx context.Context, page, count int32) ([]poststore.FeedPost, error) {
	return s.PostStore.GetAll(ctx, poststore.GetAllParams{
		Offset: page * count,
		Limit:  count,
	})
}

func (s PostService) GetFeed(ctx context.Context, args model2.GetFeedParams) (model.GetFeedResponse, error) {
	// get last updated time (for feed)
	feedUpdatedAt := s.FeedStorage.GetUpdateAt(ctx)

	var resp model.GetFeedResponse
	var feedVersion = inmemory.FeedVersionNew
	var startedFetchingAt = time.Now()

	// we have following:
	// * startedFetchingAt - time when first get request was sent (with page = 0)
	// * args.StartedFetchingAtUnixNanoUTC - the same as startedFetchingAt, but will be provided as a parameter
	// * feedUpdatedAt - the last time when sync feed process was run
	if args.Page <= 0 || args.StartedFetchingAt == nil {
		args.Page = 0
	} else if args.StartedFetchingAt.Before(feedUpdatedAt) {
		feedVersion = inmemory.FeedVersionOld
		startedFetchingAt = *args.StartedFetchingAt
	} else {
		startedFetchingAt = *args.StartedFetchingAt
	}

	// ------------------------------------ posts ------------------------------------------------------

	// fetching best-ranked feed ids for a particular page
	postIds, err := s.FeedStorage.GetFeed(ctx, feedVersion, args.Page, args.Count)
	if err != nil {
		return resp, fmt.Errorf("failed to fetch ids of posts (feed). err: %v", err)
	}

	// fetch posts by ids
	feedFetched, err := s.PostStore.GetByIds(ctx, postIds)
	if err != nil {
		return resp, fmt.Errorf("failed to fetch post by ids. err: %v", err)
	}
	length := len(feedFetched)

	sortedList := poststore.PostSortContainer(feedFetched)
	sort.Sort(sortedList)
	feedFetched = sortedList

	// ------------------------------------ ads ------------------------------------------------------

	var adCount int64
	var adIds []int64
	if length >= 17 {
		adCount = 2
	} else if length >= 3 {
		adCount = 1
	}

	// fetch ad index for a user
	// adIndexForUser -> the last ad that user has seen
	adIndexForUser, err := s.AdStorage.GetUserIndex(ctx, args.Author)
	if err != nil {
		s.Log.Printf("failed to fetch user index. err: %v \n", err)
	}

	adIds, err = s.AdStorage.GetAds(ctx, adIndexForUser, adCount)
	if err != nil {
		s.Log.Printf("failed to fetch ads. err: %v \n", err)
	}

	var promotedPosts []poststore.FeedPost
	if len(adIds) > 0 {
		promotedPosts, err = s.PostStore.GetByIds(ctx, adIds)
		if err != nil {
			s.Log.Printf("failed to fetch ads from datastore. err: %v \n", err)
		}
	}

	var j, count int
	var resultFeed = make([]poststore.FeedPost, 0, len(feedFetched)+len(promotedPosts))
	for i := 0; i < len(feedFetched); i++ {
		if (i == 1 || i == 15) && !resultFeed[i-1].Nsfw && !feedFetched[i].Nsfw && j < len(promotedPosts) {
			resultFeed = append(resultFeed, promotedPosts[j])
			count += 1
			j += 1
		}

		resultFeed = append(resultFeed, feedFetched[i])
	}

	// ------------------------------------ update values ------------------------------------------------------

	// update ad index for a user
	if err = s.AdStorage.SetUserIndex(ctx, args.Author, adIndexForUser+int64(count)); err != nil {
		s.Log.Printf("failed to update ad index for a user. err: %v \n", err)
	}

	return model.GetFeedResponse{
		Version:                      feedVersion,
		Posts:                        resultFeed,
		StartedFetchingAtUnixNanoUTC: startedFetchingAt.UTC().UnixNano(),
	}, nil
}
