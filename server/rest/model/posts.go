package model

import (
	"net/url"

	"github.com/yerassyldanay/makala/model"
	"github.com/yerassyldanay/makala/pkg/regx"
)

type PostRequest struct {
	Title     string  `json:"title" description:"title of a post"`
	Author    string  `json:"author" description:"should be a random 8 character string prefixed with t2_"  example:"t2_11qnzrqv" validate:"required"`
	Link      *string `json:"link" description:"should be a valid URL"`
	Submakala string  `json:"submakala" description:"the submakala associated with this post"`
	Content   *string `json:"content" description:"in the case of a text-only post. A post cannot have both a link and content populated"`
	Score     float64 `json:"score" description:"the total score associated with the upvotes and downvotes of a post"`
	Promoted  bool    `json:"promoted" description:"indicates whether or not the post is an ad or not"`
	Nsfw      bool    `json:"nsfw" description:"indicates whether or not the post is safe for work"`
}

func (p PostRequest) Validate() error {
	if p.Link != nil && p.Content != nil {
		return model.LinkAndContentClash
	}

	if p.Link != nil {
		_, err := url.Parse(*p.Link)
		if err != nil {
			return model.InvalidURL
		}
	}

	match := regx.ValidAuthorName.MatchString(p.Author)
	if !match {
		return model.InvalidAuthorName
	}

	return nil
}
