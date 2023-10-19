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

package huaweicloud

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/noOvertimeGroup/go-filesystem/internal/errs"

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

func isValid() bool {
	return accessKeyID == "" || accessKeySecret == "" || endpoint == "" || bucketName == ""
}

func TestStorage_PutFile(t *testing.T) {
	if isValid() {
		t.Log("huawei obs configure not found...")
		return
	}

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
				create, err := os.Create("/tmp/put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("/tmp/put.txt"))
				o := &obs.DeleteObjectInput{}
				o.Bucket = bucketName
				o.Key = target[1:]
				_, err := client.DeleteObject(o)
				require.NoError(t, err)
			},
			target: "/test/put.txt",
			file: func(t *testing.T) fs.File {
				open, err := os.Open("/tmp/put.txt")
				require.NoError(t, err)
				return open
			},
		},
		{
			name: "test huawei obs storage file path error",
			before: func(t *testing.T, target string) {
				create, err := os.Create("/tmp/put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("/tmp/put.txt"))
			},
			target: "",
			file: func(t *testing.T) fs.File {
				open, err := os.Open("/tmp/put.txt")
				require.NoError(t, err)
				return open
			},
			wantErr: errs.ErrRelativePath,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// if exist err this not run...
			defer tc.after(t, tc.target)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			s := NewStorage(client)
			tc.before(t, tc.target)
			file := tc.file(t)
			// if the file is open, it needs to be closed
			defer file.Close()
			err := s.PutFile(ctx, bucketName+tc.target, file)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestStorage_GetFile(t *testing.T) {
	if isValid() {
		t.Log("huawei obs configure not found...")
		return
	}

	client, err := obs.New(accessKeyID, accessKeySecret, endpoint)
	require.NoError(t, err)

	testCase := []struct {
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		wantVal string
		wantErr error
	}{
		{
			name: "test huawei obs storage get file",
			before: func(t *testing.T, target string) {
				var bf bytes.Buffer
				bf.WriteString("the test file...")

				input := &obs.PutObjectInput{}
				input.Bucket = bucketName
				input.Key = target[1:]
				input.Body = &bf

				_, err := client.PutObject(input)
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				input := &obs.DeleteObjectInput{}
				input.Bucket = bucketName
				input.Key = target[1:]
				_, err := client.DeleteObject(input)
				require.NoError(t, err)
			},
			target:  "/test/put.txt",
			wantVal: "the test file...",
		},
		{
			name:    "test huawei obs storage get path error",
			before:  func(t *testing.T, target string) {},
			after:   func(t *testing.T, target string) {},
			target:  "",
			wantErr: errs.ErrRelativePath,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// if exist err this not run...
			defer tc.after(t, tc.target)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			tc.before(t, tc.target)
			s := NewStorage(client)
			f, err := s.GetFile(ctx, bucketName+tc.target)
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
		t.Log("huawei obs configure not found...")
		return
	}

	client, err := obs.New(accessKeyID, accessKeySecret, endpoint)
	require.NoError(t, err)

	testCase := []struct {
		name    string
		before  func(t *testing.T, target string)
		after   func(t *testing.T, target string)
		target  string
		wantVal int64
		wantErr error
	}{
		{
			name: "test huawei obs storage get file size",
			before: func(t *testing.T, target string) {
				var bf bytes.Buffer
				bf.WriteString("the test file...")

				input := &obs.PutObjectInput{}
				input.Bucket = bucketName
				input.Key = target[1:]
				input.Body = &bf

				_, err := client.PutObject(input)
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				input := &obs.DeleteObjectInput{}
				input.Bucket = bucketName
				input.Key = target[1:]
				_, err := client.DeleteObject(input)
				require.NoError(t, err)
			},
			target:  "/test/put.txt",
			wantVal: int64(len("the test file...")),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// if exist err this not run...
			defer tc.after(t, tc.target)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			tc.before(t, tc.target)
			s := NewStorage(client)
			size, err := s.Size(ctx, bucketName+tc.target)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, size)
		})
	}
}
