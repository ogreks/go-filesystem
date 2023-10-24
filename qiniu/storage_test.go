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
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/stretchr/testify/assert"
)

var (
	accessKeyID     = os.Getenv("QINIU_OSS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("QINIU_OSS_ACCESSKEY_SECRET")
	endpoint        = os.Getenv("QINIU_OSS_ENDOPOINT") // cn-east-2
	bucketName      = os.Getenv("BUCKET")
	domain          = os.Getenv("QINIU_OSS_DOMAIN")
)

func isValid() bool {
	return accessKeyID == "" || accessKeySecret == "" || bucketName == "" || endpoint == "" || domain == ""
}

// getBucketManager get *storage.BucketManager
func getBucketManager(t *testing.T) *storage.BucketManager {
	mac := qbox.NewMac(accessKeyID, accessKeySecret)
	cfg := storage.Config{}
	// 空间对应的机房
	region, ok := storage.GetRegionByID(storage.RegionID(endpoint))
	assert.Equal(t, true, ok)
	cfg.Region = &region
	return storage.NewBucketManager(mac, &cfg)
}

func TestStorage_PutFile(t *testing.T) {
	if isValid() {
		if accessKeyID == "" {
			t.Log("accessKeyID")
		}

		if accessKeySecret == "" {
			t.Log("accessKeySecret")
		}

		if bucketName == "" {
			t.Log("bucketName")
		}

		if domain == "" {
			t.Log("domain")
		}
		t.Log("qiniu kodo configure not found...")
		return
	}

	client := getBucketManager(t)

	testCase := []struct {
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		file    io.Reader
		wantErr error
	}{
		{
			name:   "test qiniu storage put file",
			before: func(t *testing.T, target string) {},
			after: func(t *testing.T, target string) {
				require.NoError(t, client.Delete(bucketName, target))
			},
			target: "test/put.txt",
			file:   bytes.NewReader([]byte("the test file...")),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.after(t, tc.target)
			// if exist err this not run...
			tc.before(t, tc.target)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			s := NewStorage(client, domain)
			err := s.PutFile(ctx, fmt.Sprintf("%s/%s", bucketName, tc.target), tc.file)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestStorage_GetFile(t *testing.T) {
	if isValid() {
		t.Log("qiniu kodo configure not found...")
		return
	}

	client := getBucketManager(t)

	testCase := []struct {
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		wantVal string
		wantErr error
	}{
		{
			name: "test qiniu storage put file",
			before: func(t *testing.T, target string) {
				bf := bytes.NewReader([]byte("the test file..."))

				putPolicy := storage.PutPolicy{
					Scope: bucketName + ":" + target,
				}
				uploadToken := putPolicy.UploadToken(client.Mac)
				from := storage.NewFormUploader(client.Cfg)
				err := from.Put(context.Background(), nil, uploadToken, target, bf, bf.Size(), nil)
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, client.Delete(bucketName, target))
			},
			target:  "test/put.txt",
			wantVal: "the test file...",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.after(t, tc.target)
			// if exist err this not run...
			tc.before(t, tc.target)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			s := NewStorage(client, domain)

			f, err := s.GetFile(ctx, fmt.Sprintf("%s/%s", bucketName, tc.target))
			assert.Equal(t, tc.wantErr, err)
			if f == nil {
				return
			}

			content, err := io.ReadAll(f)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantVal, string(content))
		})
	}
}

func TestStorage_Size(t *testing.T) {
	if isValid() {
		t.Log("qiniu kodo configure not found...")
		return
	}

	client := getBucketManager(t)

	testCase := []struct {
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		wantVal int64
		wantErr error
	}{
		{
			name: "test qiniu storage get file size",
			before: func(t *testing.T, target string) {
				bf := bytes.NewReader([]byte("the test file..."))

				putPolicy := storage.PutPolicy{
					Scope: bucketName + ":" + target,
				}
				uploadToken := putPolicy.UploadToken(client.Mac)
				from := storage.NewFormUploader(client.Cfg)
				err := from.Put(context.Background(), nil, uploadToken, target, bf, bf.Size(), nil)
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, client.Delete(bucketName, target))
			},
			target:  "test/put.txt",
			wantVal: int64(len("the test file...")),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.after(t, tc.target)
			// if exist err this not run...
			tc.before(t, tc.target)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			s := NewStorage(client, domain)

			size, err := s.Size(ctx, fmt.Sprintf("%s/%s", bucketName, tc.target))
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, size)
		})
	}
}
