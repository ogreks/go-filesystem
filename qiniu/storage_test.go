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

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

var (
	accessKeyID     = os.Getenv("QINIU_OSS_ACCESSKEY_ID")
	accessKeySecret = os.Getenv("QINIU_OSS_ACCESSKEY_SECRET")
)

type S Client

func TestStorage_PutFile2(t *testing.T) {

	if accessKeyID == "" || accessKeySecret == "" {
		t.Log("accessKeyID/accessKeySecret not found...")
		return
	}
	bucketFile := "go-file-system"
	//上传凭据
	putPolicy := storage.PutPolicy{
		Scope: bucketFile,
	}
	mac := qbox.NewMac(accessKeyID, accessKeySecret)
	upToken := putPolicy.UploadToken(mac)

	if upToken == "" {
		t.Log("Upload Token  is nil...")
		return
	}

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	NewFormUploader := storage.NewFormUploader(&cfg)
	NewBucketManager := storage.NewBucketManager(mac, &cfg)

	s := NewStorage((*Client)(&S{
		NewFormUploader,
		NewBucketManager,
		upToken,
	}))

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
				create, err := os.Create("test_put.txt")
				require.NoError(t, err)
				defer create.Close()
				_, err = create.WriteString("the test file...")
				require.NoError(t, err)
			},
			after: func(t *testing.T, target string) {
				require.NoError(t, os.Remove("test_put.txt"))
				require.NoError(t, NewBucketManager.Delete(bucketFile, target))
			},
			target: "test_put.txt",
			file: func(t *testing.T) fs.File {
				open, err := os.Open("test_put.txt")
				require.NoError(t, err)
				return open
			},
		},
	}

	for _, tc := range testCase {

		t.Run(tc.name, func(t *testing.T) {
			//new上传类
			ctx := context.TODO()
			tc.before(t, tc.target)
			file := tc.file(t)
			err := s.PutFile(ctx, tc.target, file)
			assert.Equal(t, tc.wantErr, err)
			file.Close()
			tc.after(t, tc.target)
		})
	}
}

//
//func TestStorage_PutFile(t *testing.T) {
//
//	type fields struct {
//		client  *storage.FormUploader
//		upToken string
//	}
//	type args struct {
//		ctx    context.Context
//		target string
//		file   string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//		{
//			name: "test1",
//			fields: struct {
//				client  *storage.FormUploader
//				upToken string
//			}{client: StorageFormUploader(), upToken: getToken()},
//			args: struct {
//				ctx    context.Context
//				target string
//				file   string
//			}{ctx: context.Background(), target: "github-7.png", file: "../resources/github.png"},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Storage{
//				client:  tt.fields.client,
//				upToken: tt.fields.upToken,
//			}
//			if err := s.PutFile(tt.args.ctx, tt.args.target, tt.args.file); (err != nil) != tt.wantErr {
//				t.Errorf("PutFile() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func getToken() string {
//
//	putPolicy := storage.PutPolicy{
//		Scope: bucket,
//	}
//	mac := qbox.NewMac(accessKeyID, accessKeySecret)
//	return putPolicy.UploadToken(mac)
//}
//func StorageFormUploader() *storage.FormUploader {
//	cfg := storage.Config{}
//	// 空间对应的机房
//	cfg.Region = &storage.ZoneHuanan
//	// 是否使用https域名
//	cfg.UseHTTPS = true
//	// 上传是否使用CDN上传加速
//	cfg.UseCdnDomains = false
//	return storage.NewFormUploader(&cfg)
//
//}
//func Test(t *testing.T) {
//
//	localFile := "../resources/github.png"
//	key := "github-6.png"
//	putPolicy := storage.PutPolicy{
//		Scope: bucket,
//	}
//	mac := qbox.NewMac(accessKey, secretKey)
//	upToken := putPolicy.UploadToken(mac)
//	cfg := storage.Config{}
//	// 空间对应的机房
//	cfg.Region = &storage.ZoneHuanan
//	// 是否使用https域名
//	cfg.UseHTTPS = true
//	// 上传是否使用CDN上传加速
//	cfg.UseCdnDomains = false
//
//	// 构建表单上传的对象
//	formUploader := storage.NewFormUploader(&cfg)
//	storage := NewFileSystem(formUploader, upToken)
//	err := storage.PutFile(context.Background(), key, localFile)
//	if err != nil {
//		fmt.Print("上传失败：", err)
//	}
//	fmt.Print("上传成功！")
//}
