package adstore

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yerassyldanay/makala/datastore/redis_store"
	"github.com/yerassyldanay/makala/pkg/configx"
	"github.com/yerassyldanay/makala/provider/adstore"
)

func getAdStorage(t *testing.T) adstore.Storage {
	conf, err := configx.NewConfiguration()
	require.NoError(t, err)

	client, err := redis_store.NewRedisConnection(redis_store.RedisConnectionParams{
		Host:          conf.RedisHost,
		Port:          conf.RedisPort,
		LogicDatabase: conf.RedisAdLogicDatabase + 5,
	})
	require.NoError(t, err)

	return adstore.Storage{
		Log:    log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile),
		Client: client,
	}
}

func TestAdStore(t *testing.T) {
	st := getAdStorage(t)

	statusCmd := st.Client.FlushAll(context.Background())
	require.NotNil(t, statusCmd)
	require.NoError(t, statusCmd.Err())

	{
		// testing set/get update time
		// index must always be 0 as there are not ads
		author := "t2_author12"
		err := st.SetUserIndex(context.Background(), author, 1)
		require.NoError(t, err)

		index, err := st.GetUserIndex(context.Background(), author)
		require.NoErrorf(t, err, "failed to get update time from cache")
		require.Zerof(t, index, "updatedAtNano must be zero as there are no ads")
	}

	{
		// add ads to cache
		var i int64
		for i = 1; i <= 10; i++ {
			err := st.CreateAd(context.Background(), i)
			require.NoErrorf(t, err, "failed to create an ad in cache")
		}

		// count ads
		length, err := st.CountAds(context.Background())
		require.NoErrorf(t, err, "failed to count number of ads in cache")
		require.Equalf(t, length, int64(10), "number of ads in cache must be 10")
	}

	{
		// testing set/get update time
		author := "t2_author12"
		err := st.SetUserIndex(context.Background(), author, 1)
		require.NoError(t, err)

		index, err := st.GetUserIndex(context.Background(), author)
		require.NoErrorf(t, err, "failed to get update time from cache")
		require.NotZerof(t, index, "updatedAtNano must be non-zero as there are no ads")
	}

	{
		// testing set/get update time
		ids, err := st.GetAds(context.Background(), 0, 20)
		require.NoErrorf(t, err, "failed to fetch ads")
		require.Equalf(t, ids, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "expected different type of ads")

		// testing set/get update time
		ids, err = st.GetAds(context.Background(), 11, 20)
		require.NoErrorf(t, err, "failed to fetch ads")
		require.Equalf(t, ids, []int64{}, "expected not to get any ad")
	}
}
