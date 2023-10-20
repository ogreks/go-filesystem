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
	"io"
	"os"
	"testing"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	accessKeyID     = os.Getenv("QINIU_OSS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("QINIU_OSS_ACCESSKEY_SECRET")
	endpoint        = os.Getenv("QINIU_OSS_ENDOPOINT") // cn-east-2
	bucketName      = os.Getenv("BUCKET")
	ossDomain       = os.Getenv("QINIU_OSS_DOMAIN")
)

func TestStorage_PutFile(t *testing.T) {
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
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		file    func(t *testing.T) io.Reader
		wantErr error
	}{
		{
			name: "test qiniu storage put file",
			before: func(t *testing.T, target string) {

			},
			after: func(t *testing.T, target string) {
				//require.NoError(t, NewBucketManager.Delete(bucketName, target))
			},
			target: "test/put.txt",
			file: func(t *testing.T) io.Reader {
				bf := bytes.NewReader([]byte("the test file..."))
				return bf
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
			err := s.PutFile(ctx, tc.target, file)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestStorage_GetFile(t *testing.T) {

	if accessKeyID == "" || accessKeySecret == "" || bucketName == "" || endpoint == "" || ossDomain == "" {
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
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		wantErr error
	}{
		{
			name: "test qiniu storage put file",
			before: func(t *testing.T, target string) {

			},
			after: func(t *testing.T, target string) {

			},
			target: "test/put.txt",
		},
	}

	type MyString string
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			//new查询类
			ctx := context.WithValue(context.Background(), MyString("ossDomain"), ossDomain)
			bt, err := s.GetFile(ctx, tc.target)
			t.Log(bt)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

// TODO 注意需要等待实现
func TestStorage_Size(t *testing.T) {
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

	// TODO 这个 client 是否可以交由 具体实现去控制
	c := &Client{
		NewFormUploader,
		NewBucketManager,
		upToken,
	}

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

				err := NewFormUploader.Put(context.Background(), nil, upToken, target, bf, int64(bf.Len()), nil)
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				//require.NoError(t, NewBucketManager.Delete(bucketName, target))
			},
			target:  "test/put.txt",
			wantVal: int64(len("the test file...")),
		},
	}

	type myString string
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// if exist err this not run...
			defer tc.after(t, tc.target)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			ctx = context.WithValue(ctx, myString("bucketName"), bucketName)
			defer cancel()
			tc.before(t, tc.target)
			s := NewStorage(c)
			size, err := s.Size(ctx, tc.target)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, size)
		})
	}
}
