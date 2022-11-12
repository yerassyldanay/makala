package feedstore

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yerassyldanay/makala/datastore/redis_store"
	"github.com/yerassyldanay/makala/pkg/configx"
	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
	"github.com/yerassyldanay/makala/provider/feedstore"
)

func getFeedStorage(t *testing.T) feedstore.Storage {
	conf, err := configx.NewConfiguration()
	require.NoError(t, err)

	client, err := redis_store.NewRedisConnection(redis_store.RedisConnectionParams{
		Host:          conf.RedisHost,
		Port:          conf.RedisPort,
		LogicDatabase: conf.RedisFeedLogicDatabase + 5,
	})
	require.NoError(t, err)

	return feedstore.Storage{
		Log:    log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile),
		Client: client,
	}
}

func TestFeedStore(t *testing.T) {
	st := getFeedStorage(t)

	statusCmd := st.Client.FlushAll(context.Background())
	require.NotNil(t, statusCmd)
	require.NoError(t, statusCmd.Err())

	{
		for i := 0; i < 10; i++ {
			require.NoError(t, st.CreatePost(context.Background(), inmemory.FeedVersionTemp, rand.Float64(), rand.Int63()))
		}
	}

	{
		postIds, err := st.GetFeed(context.Background(), inmemory.FeedVersionTemp, 0, 100)
		require.NoError(t, err)
		require.Equal(t, 10, len(postIds))
	}

	{
		// remove half of the list
		require.NoError(t, st.RemovePostsWithMinScore(context.Background(), inmemory.FeedVersionTemp, 5))

		// after pop temp set must have only half of all elements
		postIds, err := st.GetFeed(context.Background(), inmemory.FeedVersionTemp, 0, 100)
		require.NoError(t, err)
		require.Equal(t, 5, len(postIds))
	}

	{
		// copy stored set to another set
		require.NoError(t, st.CopyFeed(context.Background(), inmemory.FeedVersionTemp, inmemory.FeedVersionNew))

		// check feed in each sorted set
		postIdsNew, err := st.GetFeed(context.Background(), inmemory.FeedVersionNew, 0, 100)
		require.NoError(t, err)
		require.Equal(t, 5, len(postIdsNew))

		postIdsTemp, err := st.GetFeed(context.Background(), inmemory.FeedVersionTemp, 0, 100)
		require.NoError(t, err)
		require.Equal(t, 5, len(postIdsTemp))
		require.Equal(t, postIdsNew, postIdsTemp)
	}

	{
		// empty temp sorted set
		require.NoError(t, st.CopyFeed(context.Background(), inmemory.FeedVersionEmpty, inmemory.FeedVersionTemp))
		postIdsTemp, err := st.GetFeed(context.Background(), inmemory.FeedVersionTemp, 0, 100)
		require.NoError(t, err)
		require.Equal(t, 0, len(postIdsTemp))
	}
}
