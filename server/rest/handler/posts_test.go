package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yerassyldanay/makala/pkg/convx"
	"github.com/yerassyldanay/makala/provider/poststore"
	"github.com/yerassyldanay/makala/server/rest/model"
	"github.com/yerassyldanay/makala/service/postfeed/mock"
)

type MockPostServer struct {
	PostServer
	MockPostService *mock_postfeed.MockPoster
}

func getServer(t *testing.T) *MockPostServer {
	ctr := gomock.NewController(t)
	mockServer := mock_postfeed.NewMockPoster(ctr)

	postServer := &MockPostServer{
		PostServer: PostServer{
			Router:      gin.Default(),
			PostService: mockServer,
		},
		MockPostService: mockServer,
	}
	postServer.SetRouter()
	return postServer
}

func Copy(t *testing.T, fromThis, toThis interface{}) {
	require.NoError(t, convx.Copy(fromThis, toThis))
}

func getBuffer(t *testing.T, val interface{}) io.Reader {
	b, err := json.Marshal(val)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

func TestPostServer_CreatePost(t *testing.T) {
	var payload = model.PostRequest{
		Title:     "title",
		Author:    "t2_author12",
		Link:      convx.StrToPtr("https://makala.com"),
		Submakala: "",
		Content:   nil,
		Score:     0.35,
		Promoted:  false,
		Nsfw:      false,
	}

	testCases := []struct {
		name      string
		prepare   func() (*httptest.ResponseRecorder, *http.Request)
		getServer func() *MockPostServer
		status    int
		check     func(r *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			prepare: func() (*httptest.ResponseRecorder, *http.Request) {
				req := httptest.NewRequest(http.MethodPost, "/api/makala/v1/post", getBuffer(t, payload))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				return rec, req
			},
			getServer: func() *MockPostServer {
				mockServer := getServer(t)

				var args poststore.CreateParams
				Copy(t, payload, &args)
				var argsFeedPost poststore.FeedPost
				Copy(t, payload, &argsFeedPost)

				mockServer.MockPostService.EXPECT().CreatePost(gomock.Any(), args).Return(argsFeedPost, nil)
				return mockServer
			},
			status: 200,
			check: func(rec *httptest.ResponseRecorder) {
				var result model.PostRequest
				Copy(t, payload, &result)

				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
				require.Equal(t, payload, result)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := tt.getServer()
			rec, req := tt.prepare()
			mockServer.Router.ServeHTTP(rec, req)
			assert.Equal(t, tt.status, rec.Code)
			tt.check(rec)
		})
	}
}
