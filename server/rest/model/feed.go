package model

import (
	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
	"github.com/yerassyldanay/makala/provider/poststore"
)

type GetFeedRequest struct {
	Page                         int64  `form:"page"`
	Count                        int64  `form:"count"`
	Author                       string `form:"author"`
	StartedFetchingAtUnixNanoUTC *int64 `form:"started_fetching_at_unix_nano_utc"`
}

type GetFeedResponse struct {
	Version                      inmemory.FeedVersion `json:"version"`
	Posts                        []poststore.FeedPost `json:"posts"`
	StartedFetchingAtUnixNanoUTC int64                `json:"started_fetching_at_unix_nano_utc"`
}
