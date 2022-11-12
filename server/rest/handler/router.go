package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yerassyldanay/makala/server/rest/middleware"
	"github.com/yerassyldanay/makala/service/postfeed"
)

type PostServer struct {
	Router      *gin.Engine
	PostService postfeed.Poster
}

func NewPostServer(postService postfeed.Poster) *PostServer {
	return &PostServer{
		PostService: postService,
	}
}

type RouterOption func(r *gin.Engine)

// @title           makala Feed Service
// @version         1.0.0
// @description     service stores & provides feed for users

// @contact.name   Yerassyl Danay

// @BasePath  /api/verification
func (s *PostServer) SetRouter(opts ...RouterOption) {
	router := gin.Default()

	for _, opt := range opts {
		opt(router)
	}

	router.Use(middleware.ValidateHeader())

	v1 := router.Group("/api/makala/v1/")
	{
		v1.GET("/feed", s.GetFeed)
		v1.POST("/post", s.CreatePost)
	}

	s.Router = router
}
