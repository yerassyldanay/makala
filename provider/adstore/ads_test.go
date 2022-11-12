package adstore

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_redis "github.com/yerassyldanay/makala/datastore/redis_store/mock"
	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
)

type MockStorage struct {
	Storage    Storage
	ClientMock *mock_redis.MockCmdable
}

func getMockStorage(t *testing.T) MockStorage {
	mockedClient := mock_redis.NewMockCmdable(gomock.NewController(t))

	return MockStorage{
		Storage: Storage{
			Log:    log.New(os.Stdout, "", log.LstdFlags),
			Client: mockedClient,
		},
		ClientMock: mockedClient,
	}
}

func TestStorage_SetUserIndex(t *testing.T) {
	mockStorage := getMockStorage(t)

	testCases := []struct {
		name    string
		prepare func()
		args    func() (string, int64)
		err     error
	}{
		{
			name: "ok",
			prepare: func() {
				intCmd := redis.NewIntCmd(context.Background())
				intCmd.SetVal(10)
				mockStorage.ClientMock.EXPECT().LLen(gomock.Any(), inmemory.AdsKey).Return(intCmd)

				statusCmd := redis.NewStatusCmd(context.Background())
				statusCmd.SetVal("status")
				mockStorage.ClientMock.EXPECT().Set(gomock.Any(), inmemory.UserAds("t2_author78").GetKey(), int64(1), time.Duration(0)).Return(statusCmd)
			},
			args: func() (string, int64) {
				return "t2_author78", 1
			},
			err: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			author, index := tt.args()
			err := mockStorage.Storage.SetUserIndex(context.Background(), author, index)
			require.Equal(t, tt.err, err)
		})
	}
}
