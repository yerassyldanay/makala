package adstore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"

	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
)

// Storage
// is an instance of Storer interface
// It provides all methods to work with ads
type Storage struct {
	Log    *log.Logger
	Client redis.Cmdable
}

var _ Storer = Storage{}

func (s Storage) GetAds(ctx context.Context, offset, count int64) ([]int64, error) {
	if count == 0 {
		return []int64{}, nil
	}
	cmd := s.Client.LRange(ctx, inmemory.AdsKey, offset, offset+count)
	if cmd == nil {
		return []int64{}, errors.New("empty response")
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	elements, err := cmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to parse result. err: %v", err)
	}

	var adIds = make([]int64, 0, count)
	for _, element := range elements {
		n, err := strconv.ParseInt(element, 10, 64)
		if err != nil {
			s.Log.Printf("[REDIS][GET-ADS] parse value err. err: %v \n", err)
		}
		adIds = append(adIds, n)
	}

	return adIds, nil
}

func (s Storage) CreateAd(ctx context.Context, postId int64) error {
	intCmd := s.Client.RPush(ctx, inmemory.AdsKey, postId)
	if intCmd == nil {
		return fmt.Errorf("failed to add an ad. err: empty response")
	} else if intCmd.Err() != nil {
		return fmt.Errorf("failed to add an add. err: %v", intCmd.Err())
	}
	return nil
}

func (s Storage) adsCount(ctx context.Context) (int64, error) {
	intCmd := s.Client.LLen(ctx, inmemory.AdsKey)
	if intCmd == nil {
		return 0, fmt.Errorf("failed to count elements. err: empty response")
	}

	res, err := intCmd.Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count ads. err: %v", err)
	}

	return res, nil
}

func (s Storage) CountAds(ctx context.Context) (int64, error) {
	return s.adsCount(ctx)
}

func (s Storage) SetUserIndex(ctx context.Context, author string, index int64) error {
	length, err := s.adsCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get ads count. err: %v", err)
	}
	if length > 0 {
		index = index % length
	} else {
		index = 0
	}

	statusCmd := s.Client.Set(ctx, inmemory.UserAds(author).GetKey(), index, 0)
	if statusCmd != nil && statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func (s Storage) GetUserIndex(ctx context.Context, author string) (int64, error) {
	stringCmd := s.Client.Get(ctx, inmemory.UserAds(author).GetKey())
	if stringCmd == nil {
		return 0, fmt.Errorf("failed to get ad index for user. err: empty response")
	}
	if stringCmd != nil && stringCmd.Err() != nil {
		return 0, stringCmd.Err()
	}

	numStr, err := stringCmd.Result()
	if err != nil {
		return 0, err
	}

	n, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse index. err: %v", err)
	}

	length, err := s.adsCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get length. err: %v", err)
	}

	if length == 0 {
		n = 0
	} else {
		n = n % length
	}

	return n, err
}
