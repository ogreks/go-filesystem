package aliyun

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

var (
	accessKeyID     = os.Getenv("ALIYUN_OSS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("ALIYUN_OSS_ACCESSKEY_SECRET")
	endpoint        = os.Getenv("ALIYUN_OSS_ENDPOINT")
	bucketName      = os.Getenv("ALIYUN_OSS_BUCKET")
)

func TestStorage_PutFile(t *testing.T) {
	assert.NotEmpty(t, accessKeyID)
	assert.NotEmpty(t, accessKeySecret)
	assert.NotEmpty(t, endpoint)
	assert.NotEmpty(t, bucketName)

	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	require.NoError(t, err)

	bucket, err := client.Bucket(bucketName)
	require.NoError(t, err)

	testCase := []struct {
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		file    func(t *testing.T) fs.File
		wantErr error
	}{
		{
			name: "test storage put file",
			before: func(t *testing.T, target string) {
				create, err := os.Create("put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("put.txt"))
				require.NoError(t, bucket.DeleteObject(target))
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
			ctx := context.TODO()
			s := NewStorage(bucket)
			tc.before(t, tc.target)
			file := tc.file(t)
			err := s.PutFile(ctx, tc.target, file)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t, tc.target)
			file.Close()
		})
	}
}

func TestStorage_GetFile(t *testing.T) {

}