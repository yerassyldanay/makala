package regx

import "regexp"

var (
	ValidAuthorName = regexp.MustCompile("^t2_[a-zA-Z0-9]{8}$")
	TextOnly        = regexp.MustCompile("[a-zA-Z|[[:space:]]")
)
