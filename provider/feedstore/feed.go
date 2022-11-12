package feedstore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/yerassyldanay/makala/pkg/configx"
	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
)

type Storage struct {
	Conf   configx.Configuration
	Log    *log.Logger
	Client redis.Cmdable
}

var _ Storer = Storage{}

func (s Storage) CreatePost(ctx context.Context, version inmemory.FeedVersion, score float64, id int64) error {
	cmd := s.Client.ZAdd(ctx, version.GetKey(), &redis.Z{
		Score:  score,
		Member: id,
	})
	if cmd == nil {
		return nil
	}
	return cmd.Err()
}

func (s Storage) RemovePostsWithMinScore(ctx context.Context, feedVersion inmemory.FeedVersion, maxElements int64) error {
	intCmd := s.Client.ZCount(ctx, feedVersion.GetKey(), "-inf", "+inf")
	if intCmd == nil {
		return fmt.Errorf("failed to count the elements of feed %s. err: nil val", feedVersion.GetKey())
	}

	n, err := intCmd.Result()
	if err != nil {
		return fmt.Errorf("failed to fetch redis result. err: %v", err)
	}

	for i := n - maxElements; i > 0; i-- {
		intSliceCmd := s.Client.ZPopMin(ctx, feedVersion.GetKey())
		if intSliceCmd == nil {
			s.Log.Println("[REDIS-POP] failed to pop element with min score (pop-min)")
		} else if intSliceCmd.Err() != nil {
			s.Log.Printf("[REDIS-POP] failed to pop a post with min score. err: %v \n", err)
		}
	}

	return nil
}

func (s Storage) CopyFeed(ctx context.Context, fromThis, toThis inmemory.FeedVersion) error {
	intCmd := s.Client.ZRangeStore(ctx, toThis.GetKey(), redis.ZRangeArgs{
		Key:     fromThis.GetKey(),
		Start:   "-inf",
		Stop:    "+inf",
		ByScore: true,
	})
	if intCmd == nil {
		return fmt.Errorf("received an empty response while copying sorted sets")
	} else if intCmd.Err() != nil {
		return fmt.Errorf("failed to copy sorted set. err: %v", intCmd.Err())
	}
	return nil
}

func (s Storage) GetFeed(ctx context.Context, feedVersion inmemory.FeedVersion, page, count int64) ([]int64, error) {
	stringSliceCmd := s.Client.ZRevRangeByScore(ctx, feedVersion.GetKey(), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: page * count,
		Count:  count,
	})
	if stringSliceCmd == nil && stringSliceCmd.Err() != nil {
		return nil, fmt.Errorf("failed to fetch feed from sorted set. err: %v", stringSliceCmd.Err())
	}

	result, err := stringSliceCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed from sorted set. err: %v", err)
	}

	var postIds = make([]int64, 0, count)
	for _, eachId := range result {
		num, err := strconv.ParseInt(eachId, 10, 64)
		if err != nil {
			s.Log.Printf("failed to parse post_id in redis (fetch). err: %v", err)
			continue
		}
		postIds = append(postIds, num)
	}
	return postIds, nil
}

func (s Storage) GetFeedLength(ctx context.Context, feedVersion inmemory.FeedVersion) (int64, error) {
	intCmd := s.Client.ZCount(ctx, feedVersion.GetKey(), "-inf", "+inf")
	if intCmd == nil {
		return 0, errors.New("failed to count the number of element of sorted set")
	}

	if err := intCmd.Err(); err != nil {
		return 0, fmt.Errorf("failed to count the number of element of sorted set. err: %v", err)
	}

	return intCmd.Result()
}

func (s Storage) SetUpdateAt(ctx context.Context, updatedAt time.Time) error {
	cmd := s.Client.Set(ctx, inmemory.LastUpdatedAtKey.GetKey(), updatedAt.Local().UnixNano(), -1)
	if cmd == nil {
		return errors.New("empty response")
	}
	if err := cmd.Err(); err != nil {
		return fmt.Errorf("failed to set updated time. err: %v", err)
	}

	return nil
}

func (s Storage) GetUpdateAt(ctx context.Context) time.Time {
	now := time.Now()
	cmd := s.Client.Get(ctx, inmemory.LastUpdatedAtKey.GetKey())
	if cmd == nil {
		s.Log.Printf("[FEED-STORE] empty response \n")
		return now
	}
	if err := cmd.Err(); err != nil {
		return now
	}

	updatedAt, err := cmd.Int64()
	if err != nil {
		s.Log.Printf("[FEED-STORE] failed to fetch updated time. err: %v", err)
		return now
	}

	return time.Unix(0, updatedAt)
}
