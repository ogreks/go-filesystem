package tencentyun

import (
	"context"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	accessKeyID     = os.Getenv("TENCENT_OSS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("TENCENT_OSS_ACCESSKEY_SECRET")
	endpoint        = os.Getenv("TENCENT_OSS_ENDPOINT")
)

func TestStorage_PutFile(t *testing.T) {
	if accessKeyID == "" {
		t.Log("tencent oss configure not found...")
		return
	}
	assert.NotEmpty(t, accessKeyID)
	assert.NotEmpty(t, accessKeySecret)
	assert.NotEmpty(t, endpoint)

	endpointUrl, err := url.Parse(endpoint)
	require.NoError(t, err)

	b := &cos.BaseURL{BucketURL: endpointUrl}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  accessKeyID,
			SecretKey: accessKeySecret,
		},
	})

	testCase := []struct {
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		file    func(t *testing.T) fs.File
		wantErr error
	}{
		{
			name: "test tencent oss storage put file",
			before: func(t *testing.T, target string) {
				create, err := os.Create("put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("put.txt"))
				_, err := client.Object.Delete(context.Background(), target)
				require.NoError(t, err)
			},
			target: "test/put.txt",
			file: func(t *testing.T) fs.File {
				open, err := os.Open("put.txt")
				require.NoError(t, err)
				return open
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := NewStorage(client)
			tc.before(t, tc.target)
			file := tc.file(t)
			err := s.PutFile(ctx, tc.target, file)
			assert.Equal(t, tc.wantErr, err)
			file.Close()
			tc.after(t, tc.target)
		})
	}
}

func TestStorage_GetFile(t *testing.T) {

}
