package yandexDisk

import (
	"context"
	"github.com/nikitaksv/yandex-disk-sdk-go"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func getYandexDisk() *yandexDisk {
	yadClient, err := yadisk.NewYaDisk(context.Background(), http.DefaultClient, &yadisk.Token{AccessToken: os.Getenv("YANDEX_TOKEN")})
	if err != nil {
		panic(err)
	}

	return &yandexDisk{client: yadClient}
}

func TestNewYandexDisk(t *testing.T) {

	type args struct {
		ctx    context.Context
		client *http.Client
		token  string
	}
	tests := []struct {
		name string
		args args
		want *yandexDisk
	}{
		{"success", args{context.Background(), http.DefaultClient, os.Getenv("YANDEX_TOKEN")}, getYandexDisk()},
		{"panic", args{context.Background(), http.DefaultClient, ""}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				if got := NewYandexDisk(tt.args.ctx, tt.args.client, tt.args.token); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewYandexDisk() = %v, want %v", got, tt.want)
				}
			} else {
				defer func() {
					if r := recover().(error); r != nil {
						if r.Error() != "required token" {
							t.Errorf("NewYandexDisk() panic = %v", r)
						}
					}
				}()
				_ = NewYandexDisk(tt.args.ctx, tt.args.client, tt.args.token)
			}

		})
	}
}

func Test_yandexDisk_Info(t *testing.T) {
	type fields struct {
		yd *yandexDisk
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"success", fields{getYandexDisk()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("yandexDisk.Info() panic = %v", r)
				}
			}()
			got, err := tt.fields.yd.Info()
			if err != nil {
				t.Errorf("yandexDisk.Info() error = %v", err)
			}
			if got == new(diskInfo) {
				t.Errorf("yandexDisk.Info() empty info")
			}
		})
	}
}

func Test_yandexDisk_getResourceInfo(t *testing.T) {
	type fields struct {
		yd *yandexDisk
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"file", fields{getYandexDisk()}, args{"/test/forTest.docx"}, false},
		{"directory", fields{getYandexDisk()}, args{"/test"}, false},
		{"error", fields{getYandexDisk()}, args{"/nonexistentFile.jpg"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			yd := tt.fields.yd
			got, err := yd.getResourceInfo(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("yandexDisk.getResourceInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == new(resourceInfo) {
				t.Errorf("yandexDisk.getResourceInfo() empty info")
			}
			if tt.name == "directory" {
				res := got.Resources()
				for _, rRes := range res {
					if rRes.IsDir() && rRes.Name() == "sub" {
						resD := rRes.Resources()
						find := false
						for _, rResd := range resD {
							if rResd.Name() == "subForTest.docx" && !rResd.IsDir() {
								find = true
							}
							if rResd.ParentResource() == nil {
								t.Errorf("yandexDisk.getResourceInfo() parent resource not exists")
							}
						}
						if !find {
							t.Errorf("yandexDisk.getResourceInfo() subResource subForFile.docx not exists")
						}
					}
				}
			}
		})
	}
}

func Test_convertPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"/", args{"/test/1.png"}, "disk:/test/1.png"},
		{"without /", args{"test/1.png"}, "disk:/test/1.png"},
		{"disk:/", args{"disk:/test/1.png"}, "disk:/test/1.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertPath(tt.args.path); got != tt.want {
				t.Errorf("convertPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
