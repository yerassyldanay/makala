package model

import (
	"time"

	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
	"github.com/yerassyldanay/makala/provider/poststore"
)

type CreatePostParamsExtended struct {
	CreateParams poststore.CreateParams
	FeedVersion  inmemory.FeedVersion
}

type GetFeedParams struct {
	Page              int64      `json:"page"`
	Count             int64      `json:"count"`
	Author            string     `json:"author"`
	StartedFetchingAt *time.Time `json:"started_fetching_at"`
}
