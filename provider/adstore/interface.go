package adstore

import "context"

// Storer abstracts method for working with ads (promoted posts)
// * getting a list of ads
// * creating ads
// * getting & updating user index
type Storer interface {
	// GetAds returns 'count' number of ads starting from 'offset'
	// order by the time of addition
	GetAds(ctx context.Context, offset, count int64) ([]int64, error)
	// CreateAd adds an add to the list
	CreateAd(ctx context.Context, postId int64) error
	// CountAds returns length of ads in the list
	CountAds(ctx context.Context) (int64, error)
	// SetUserIndex sets index for each user that users will be given an ad where it stopped
	SetUserIndex(ctx context.Context, author string, index int64) error
	// GetUserIndex returns an ad index for a given user
	GetUserIndex(ctx context.Context, author string) (int64, error)
}
