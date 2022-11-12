package postfeed

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/yerassyldanay/makala/pkg/configx"
	"github.com/yerassyldanay/makala/pkg/constx/inmemory"
	"github.com/yerassyldanay/makala/pkg/convx"
	mock_adstore "github.com/yerassyldanay/makala/provider/adstore/mock"
	mock_feedstore "github.com/yerassyldanay/makala/provider/feedstore/mock"
	"github.com/yerassyldanay/makala/provider/poststore"
	mock_poststore "github.com/yerassyldanay/makala/provider/poststore/mock"
)

type MockPostService struct {
	Service     PostService
	PostStore   *mock_poststore.MockQuerier
	FeedStorage *mock_feedstore.MockStorer
	AdStorage   *mock_adstore.MockStorer
}

func NewMockPostService(t *testing.T) MockPostService {
	ctrl := gomock.NewController(t)

	mockedPostStore := mock_poststore.NewMockQuerier(ctrl)
	mockedFeedStorage := mock_feedstore.NewMockStorer(ctrl)
	mockedAdStorage := mock_adstore.NewMockStorer(ctrl)

	return MockPostService{
		Service: PostService{
			Conf:        configx.Configuration{},
			Log:         log.New(os.Stdout, "", log.LstdFlags),
			PostStore:   mockedPostStore,
			FeedStorage: mockedFeedStorage,
			AdStorage:   mockedAdStorage,
		},
		PostStore:   mockedPostStore,
		FeedStorage: mockedFeedStorage,
		AdStorage:   mockedAdStorage,
	}
}

func TestPostService_CreatePost(t *testing.T) {
	mockService := NewMockPostService(t)

	postOrdinary := poststore.FeedPost{
		ID:        1,
		Title:     "title",
		Author:    "t2_author23",
		Link:      convx.StrToPtr("https://makala.com"),
		Submakala: "submakala fav",
		Score:     3.35,
		Promoted:  false,
	}

	postPromoted := poststore.FeedPost{
		ID:        2,
		Title:     "title 2",
		Author:    "t2_author34",
		Link:      convx.StrToPtr("https://makala.com"),
		Submakala: "submakala fav",
		Score:     3.12,
		Promoted:  true,
	}

	testCases := []struct {
		name    string
		args    func() poststore.CreateParams
		resp    func() poststore.FeedPost
		prepare func()
		err     error
	}{
		{
			name: "ok",
			args: func() poststore.CreateParams {
				var args poststore.CreateParams
				require.NoError(t, convx.Copy(postOrdinary, &args))
				return args
			},
			resp: func() poststore.FeedPost {
				return postOrdinary
			},
			prepare: func() {
				var args poststore.CreateParams
				require.NoError(t, convx.Copy(postOrdinary, &args))

				mockService.PostStore.EXPECT().Create(gomock.Any(), args).
					Return(postOrdinary, nil)
				mockService.FeedStorage.EXPECT().GetFeedLength(gomock.Any(), inmemory.FeedVersionNew).
					Return(int64(0), errors.New("escape this test error"))
				mockService.FeedStorage.EXPECT().CreatePost(gomock.Any(), inmemory.FeedVersionNew, postOrdinary.Score, postOrdinary.ID).
					Return(errors.New("escape this test error"))
			},
			err: nil,
		},
		{
			name: "error",
			args: func() poststore.CreateParams {
				var args poststore.CreateParams
				require.NoError(t, convx.Copy(postOrdinary, &args))
				return args
			},
			resp: func() poststore.FeedPost {
				return poststore.FeedPost{}
			},
			prepare: func() {
				var args poststore.CreateParams
				require.NoError(t, convx.Copy(postOrdinary, &args))

				mockService.PostStore.EXPECT().Create(gomock.Any(), args).
					Return(postOrdinary, errors.New("failed to create"))
			},
			err: errors.New("failed to create"),
		},
		{
			name: "ok-ad",
			args: func() poststore.CreateParams {
				var args poststore.CreateParams
				require.NoError(t, convx.Copy(postPromoted, &args))
				return args
			},
			resp: func() poststore.FeedPost {
				return postPromoted
			},
			prepare: func() {
				var args poststore.CreateParams
				require.NoError(t, convx.Copy(postPromoted, &args))

				mockService.PostStore.EXPECT().Create(gomock.Any(), args).
					Return(postPromoted, nil)
				mockService.FeedStorage.EXPECT().GetFeedLength(gomock.Any(), inmemory.FeedVersionNew).
					Return(int64(0), errors.New("escape this test error"))
				mockService.AdStorage.EXPECT().CreateAd(gomock.Any(), postPromoted.ID).
					Return(errors.New("escape this test error"))
			},
			err: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			resp, err := mockService.Service.CreatePost(context.Background(), tt.args())
			require.Equal(t, tt.err != nil, err != nil)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.resp(), resp)
		})
	}
}
