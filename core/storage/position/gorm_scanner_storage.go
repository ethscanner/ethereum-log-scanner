package position

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethscanner/ethereum-log-scanner/internal/dao"
	"github.com/ethscanner/ethereum-log-scanner/internal/model/entity"
)

type gOrmScannerStorage struct {
}

func NewGormScannerStorage() *gOrmScannerStorage {
	return &gOrmScannerStorage{}
}

func (s *gOrmScannerStorage) UpdateUint64(ctx context.Context, key string, val uint64) error {
	ret, err := dao.HdSystemConfig.Ctx(ctx).Data("value=?", val).Where("`key`=?", key).Update()
	if err != nil {
		return err
	} else if rows, _ := ret.RowsAffected(); rows != 1 {
		return errors.New("update db error")
	}
	return nil
}

func (s *gOrmScannerStorage) InsertUint64(ctx context.Context, key string, val uint64) error {
	_, err := dao.HdSystemConfig.Ctx(ctx).InsertIgnore(entity.HdSystemConfig{
		Name:        key,
		Key:         key,
		Value:       fmt.Sprint(val),
		Description: "",
	})
	return err
}

func (s *gOrmScannerStorage) GetUint64ByKey(ctx context.Context, key string) (val uint64, err error) {
	ret, err := dao.HdSystemConfig.Ctx(ctx).Value("value", "`key`=?", key)
	if err != nil {
		return 0, err
	}
	return ret.Uint64(), nil
}
