package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yerassyldanay/makala/pkg/convx"
	"github.com/yerassyldanay/makala/provider/poststore"
	"github.com/yerassyldanay/makala/server/rest/model"
	model2 "github.com/yerassyldanay/makala/service/postfeed/model"
)

// CreatePost
// @Tags post
// @Summary creates a post
// @Description creates a post
// @Accept  json
// @Produce  json
// @Param args body model.PostRequest true "post info"
// @Success 200 {object} poststore.FeedPost
// @Failure 400 {object} model.ErrMsg
// @Router /api/makala/v1/post [POST]
func (s PostServer) CreatePost(c *gin.Context) {
	var args model.PostRequest
	if err := c.ShouldBindJSON(&args); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrMsg{
			Err: fmt.Sprintf("failed to parse params. err: %v", err),
		})
		return
	}

	if err := args.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrMsg{
			Err: fmt.Sprintf("invalid payload. err: %v", err),
		})
		return
	}

	var post poststore.CreateParams
	if err := convx.Copy(args, &post); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrMsg{
			Err: fmt.Sprintf("failed to copy params. err: %v", err),
		})
		return
	}

	feedPost, err := s.PostService.CreatePost(context.Background(), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrMsg{
			Err: fmt.Sprintf("failed to create a post. err: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, feedPost)
}

// GetFeed
// @Tags feed
// @Summary fetches feed
// @Description this API lets a user fetch feed
// @Accept  json
// @Produce  json
// @Param        page    query     int32  false  "page"
// @Param        count    query     int32  false  "number of elements"
// @Param        author    query     string  false  "author name for ad index caching"
// @Param        started_fetching_at_unix_nano_utc    query     int64  false  "time, when a user started fetching feed"
// @Success 200 {object} model.GetFeedResponse
// @Failure 400 {object} model.ErrMsg
// @Router /api/makala/v1/feed [GET]
func (s PostServer) GetFeed(c *gin.Context) {
	var args = model.GetFeedRequest{
		Page:  0,
		Count: 27,
	}
	if err := c.ShouldBindQuery(&args); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrMsg{
			Err: fmt.Sprintf("failed to parse params. err: %v", err),
		})
		return
	}

	var startedFetchingAt *time.Time
	if args.StartedFetchingAtUnixNanoUTC != nil {
		startedFetchingAt = convx.TimeToPtr(time.Unix(0, *args.StartedFetchingAtUnixNanoUTC))
	}

	// logic for getting a list of makala feed
	getFeedResp, err := s.PostService.GetFeed(context.Background(), model2.GetFeedParams{
		Page:              args.Page,
		Count:             args.Count,
		Author:            args.Author,
		StartedFetchingAt: startedFetchingAt,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrMsg{
			Err: fmt.Sprintf("failed to fetch feed. err: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, model.GetFeedResponse{
		Version:                      getFeedResp.Version,
		StartedFetchingAtUnixNanoUTC: getFeedResp.StartedFetchingAtUnixNanoUTC,
		Posts:                        getFeedResp.Posts,
	})
}
