package huaweicloud

import (
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	accessKeyID     = os.Getenv("HUAWEI_OBS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("HUAWEI_OBS_ACCESSKEY_SECRET")
	endpoint        = os.Getenv("HUAWEI_OBS_ENDPOINT")
	bucketName      = os.Getenv("BUCKET")
)

func TestStorage_PutFile(t *testing.T) {
	if accessKeyID == "" {
		t.Log("huawei obs configure not found...")
		return
	}

	assert.NotEmpty(t, accessKeyID)
	assert.NotEmpty(t, accessKeySecret)
	assert.NotEmpty(t, endpoint)
	assert.NotEmpty(t, bucketName)

	client, err := obs.New(accessKeyID, accessKeySecret, endpoint)
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
			name: "test huawei obs storage put file",
			before: func(t *testing.T, target string) {
				create, err := os.Create("put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("put.txt"))
				o := &obs.DeleteObjectInput{}
				o.Bucket = bucketName
				o.Key = target[1:]
				_, err := client.DeleteObject(o)
				require.NoError(t, err)
			},
			target: "/test/put.txt",
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
			s := NewStorage(client)
			tc.before(t, tc.target)
			file := tc.file(t)
			err := s.PutFile(ctx, bucketName+tc.target, file)
			assert.Equal(t, tc.wantErr, err)
			tc.after(t, tc.target)
			file.Close()
		})
	}
}

func TestStorage_GetFile(t *testing.T) {
}
