// Copyright (c) 2023 noOvertimeGroup
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package qiniu

import (
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	accessKeyID     = os.Getenv("QINIU_OSS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("QINIU_OSS_ACCESSKEY_SECRET")
	endpoint        = os.Getenv("QINUI_OSS_ENDOPOINT") // cn-east-2
	bucketName      = os.Getenv("BUCKET")
)

func TestStorage_PutFile2(t *testing.T) {
	if accessKeyID == "" || accessKeySecret == "" || bucketName == "" || endpoint == "" {
		t.Log("qiniu kodo configure not found...")
		return
	}
	//上传凭据
	putPolicy := storage.PutPolicy{
		Scope: bucketName,
	}
	mac := qbox.NewMac(accessKeyID, accessKeySecret)
	upToken := putPolicy.UploadToken(mac)

	if upToken == "" {
		t.Log("Upload Token  is nil...")
		return
	}

	cfg := storage.Config{}
	// 空间对应的机房
	region, ok := storage.GetRegionByID(storage.RegionID(endpoint))
	assert.Equal(t, true, ok)
	cfg.Region = &region
	NewFormUploader := storage.NewFormUploader(&cfg)
	NewBucketManager := storage.NewBucketManager(mac, &cfg)

	s := NewStorage(&Client{
		NewFormUploader,
		NewBucketManager,
		upToken,
	})

	testCase := []struct {
		name       string
		before     func(t *testing.T, target string)
		after      func(t *testing.T, target string)
		target     string
		bucketFile string
		file       func(t *testing.T) fs.File
		wantErr    error
	}{
		{
			name: "test qiniu storage put file",
			before: func(t *testing.T, target string) {
				create, err := os.Create("/tmp/test_put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("/tmp/test_put.txt"))
				require.NoError(t, NewBucketManager.Delete(bucketName, target))
			},
			target: "test_put.txt",
			file: func(t *testing.T) fs.File {
				open, err := os.Open("/tmp/test_put.txt")
				require.NoError(t, err)
				return open
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.after(t, tc.target)
			//new上传类
			ctx := context.TODO()
			tc.before(t, tc.target)
			file := tc.file(t)
			defer file.Close()
			err := s.PutFile(ctx, tc.target, file)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
