package adstore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

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

const maxRequests = 10
const maxSuccess = 4

func makeRequest(ctx context.Context, urlString string, successChan chan struct{}) {
	if len(successChan) == cap(successChan) {
		fmt.Println("we have enough number of successful requests...")
		return
	}

	cl := http.DefaultClient
	resp, err := cl.Get(urlString)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("%q - not ok (%d). err: %v \n", urlString, resp.StatusCode, err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("%q - ok (%d) \n", urlString, resp.StatusCode)
	successChan <- struct{}{}
}

func handleUrls(ch <-chan string) {
	var wg = &sync.WaitGroup{}
	var limitCh = make(chan struct{}, maxRequests)

	var successChan = make(chan struct{}, maxSuccess)

	for urlString := range ch {

		limitCh <- struct{}{}

		wg.Add(1)
		go func(urlString string, wg *sync.WaitGroup) {
			defer wg.Done()
			makeRequest(context.Background(), urlString, successChan)
			<-limitCh
		}(urlString, wg)
	}
}

func TestUnicode(t *testing.T) {
	urls := make([]string, 0, 100)
	for i := 0; i < 50; i++ {
		urls = append(urls, []string{
			"https://google.com",
			"http://google.com",
		}...)
	}

	var urlChan = make(chan string)
	go func() {
		for _, urlString := range urls {
			urlChan <- urlString
		}
		close(urlChan)
	}()

	handleUrls(urlChan)
}

func TestGoroutine(t *testing.T) {
	ctx, can := context.WithCancel(context.Background())
	for i := 0; i < 100; i++ {
		go func(ctx context.Context, i int) {
			select {
			case <-ctx.Done():
			default:
				time.Sleep(time.Second * 1)
				fmt.Printf("finished the goroutine %d \n", i)
			}
		}(ctx, i)
	}

	time.Sleep(time.Second * 1)
	can()
}
