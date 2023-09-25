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
	if accessKeyID == "" || accessKeySecret == "" || endpoint == "" {
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

}
