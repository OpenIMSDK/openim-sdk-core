package db

import (
	"context"
	"testing"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/tools/log"
)

func Test_GetAppSDKVersion(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}
	info, err := db.GetAppSDKVersion(ctx)
	if err != nil {
		t.Fatal(err)
	}
	log.ZDebug(ctx, "info is ", "info", info)
}

func Test_SetAppSDKVersion(t *testing.T) {
	ctx := context.Background()
	db, err := NewDataBase(ctx, "1695766238", "./", 6)
	if err != nil {
		return
	}

	newVersion := &model_struct.LocalAppSDKVersion{
		Version: "4.0.0",
	}

	err = db.SetAppSDKVersion(ctx, newVersion)
	if err != nil {
		t.Fatal(err)
	}
}
