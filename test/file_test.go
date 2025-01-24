package test

import (
	"path/filepath"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
)

func TestUploadFile(t *testing.T) {

	fp := `C:\Users\openIM\Desktop\dist.zip`

	resp, err := open_im_sdk.UserForSDK.File().UploadFile(ctx, &file.UploadFileReq{
		Filepath: fp,
		Name:     filepath.Base(fp),
		Cause:    "test",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
