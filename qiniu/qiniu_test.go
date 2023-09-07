package qiniu

import (
	"context"
	"fmt"
	"github.com/noOvertimeGroup/go-filesystem"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"reflect"
	"testing"
)

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

var bucket = "go-file-system"
var accessKey = "accessKey"
var secretKey = "secretKey"

func TestNewFileSystem(t *testing.T) {

	type args struct {
		bucket  *storage.FormUploader
		upToken string
	}
	tests := []struct {
		name string
		args args
		want filesystem.FileSystem
	}{
		{
			name: "test1",
			args: args{
				bucket:  StorageFormUploader(),
				upToken: getToken(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileSystem(tt.args.bucket, tt.args.upToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_PutFile(t *testing.T) {
	type fields struct {
		client  *storage.FormUploader
		upToken string
	}
	type args struct {
		ctx    context.Context
		target string
		file   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: struct {
				client  *storage.FormUploader
				upToken string
			}{client: StorageFormUploader(), upToken: getToken()},
			args: struct {
				ctx    context.Context
				target string
				file   string
			}{ctx: context.Background(), target: "github-6.png", file: "../resources/github.png"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				client:  tt.fields.client,
				upToken: tt.fields.upToken,
			}
			if err := s.PutFile(tt.args.ctx, tt.args.target, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("PutFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getToken() string {

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	return putPolicy.UploadToken(mac)
}
func StorageFormUploader() *storage.FormUploader {
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	return storage.NewFormUploader(&cfg)

}
func Test(t *testing.T) {

	localFile := "../resources/github.png"
	key := "github-6.png"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	storage := NewFileSystem(formUploader, upToken)
	err := storage.PutFile(context.Background(), key, localFile)
	if err != nil {
		fmt.Print("上传失败：", err)
	}
	fmt.Print("上传成功！")
}
