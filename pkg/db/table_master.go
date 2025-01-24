//go:build !js
// +build !js

package db

import (
	"context"

	"github.com/openimsdk/tools/errs"
)

func (d *DataBase) GetExistTables(ctx context.Context) ([]string, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var tables []string
	return tables, errs.Wrap(d.conn.WithContext(ctx).Raw("SELECT name FROM sqlite_master WHERE type='table'").Scan(&tables).Error)
}
