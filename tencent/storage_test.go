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

package tencent

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	accessKeyID     = os.Getenv("TENCENT_OSS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("TENCENT_OSS_ACCESSKEY_SECRET")
	endpoint        = os.Getenv("TENCENT_OSS_ENDPOINT")
)

func isValid() bool {
	return accessKeyID == "" || accessKeySecret == "" || endpoint == ""
}

func getClient() (*cos.Client, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	b := &cos.BaseURL{BucketURL: endpointUrl}

	return cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  accessKeyID,
			SecretKey: accessKeySecret,
		},
	}), nil
}

func TestStorage_PutFile(t *testing.T) {
	if isValid() {
		t.Log("tencent oss configure not found...")
		return
	}

	client, err := getClient()
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
			name: "test tencent oss storage put file",
			before: func(t *testing.T, target string) {
				create, err := os.Create("/tmp/test_put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("/tmp/test_put.txt"))
				_, err := client.Object.Delete(context.Background(), target)
				require.NoError(t, err)
			},
			target: "test/put.txt",
			file: func(t *testing.T) fs.File {
				open, err := os.Open("/tmp/test_put.txt")
				require.NoError(t, err)
				return open
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// if exist err this not run...
			defer tc.after(t, tc.target)
			ctx := context.Background()
			s := NewStorage(client)
			tc.before(t, tc.target)

			file := tc.file(t)
			// if the file is open, it needs to be closed
			defer file.Close()
			err := s.PutFile(ctx, tc.target, file)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestStorage_GetFile(t *testing.T) {
	if isValid() {
		t.Log("tencent oss configure not found...")
		return
	}

	client, err := getClient()
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
			name: "test tencent storage get file",
			before: func(t *testing.T, target string) {
				var bf bytes.Buffer
				bf.WriteString("the test file...")

				_, err := client.Object.Put(context.Background(), target, &bf, nil)
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				_, err := client.Object.Delete(context.Background(), target)
				require.NoError(t, err)
			},
			target:  "test/put.txt",
			wantVal: "the test file...",
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
			f, err := s.GetFile(ctx, tc.target)
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
		t.Log("tencent oss configure not found...")
		return
	}

	client, err := getClient()
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
			name: "test tencent storage get file size",
			before: func(t *testing.T, target string) {
				var bf bytes.Buffer
				bf.WriteString("the test file...")

				_, err := client.Object.Put(context.Background(), target, &bf, nil)
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				_, err := client.Object.Delete(context.Background(), target)
				require.NoError(t, err)
			},
			target:  "test/put.txt",
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
			size, err := s.Size(ctx, tc.target)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, size)
		})
	}
}
