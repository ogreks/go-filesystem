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

package local

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestBucket_GetDir(t *testing.T) {
	type fields struct {
		filepath string
		file     *os.File
		fileInfo os.FileInfo
	}
	type args struct {
		filepath string
	}
	filepath, _ := filepath.Abs("./bucket.go")
	tt := struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		name:   "测试1",
		fields: fields{},
		args: args{
			filepath: filepath,
		},
	}
	fmt.Println(tt.want)
	t.Run(tt.name, func(t *testing.T) {
		b := &Bucket{
			filepath: tt.fields.filepath,
			file:     tt.fields.file,
			fileInfo: tt.fields.fileInfo,
		}
		if got := b.GetDir(tt.args.filepath); got != tt.want {
			t.Errorf("GetDir() = %v, want %v", got, tt.want)
		}
	})
}
