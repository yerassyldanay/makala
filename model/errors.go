package model

import "errors"

type PostRequestError error

var (
	InvalidAuthorName   PostRequestError = errors.New("invalid author name. It should be a random 8 character string prefixed with t2_")
	InvalidURL          PostRequestError = errors.New("invalid url")
	LinkAndContentClash PostRequestError = errors.New("post cannot bear content & link at the same time")
)
