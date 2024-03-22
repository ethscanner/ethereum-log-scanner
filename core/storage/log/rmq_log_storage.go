package log

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethscanner/ethereum-log-scanner/core/cache"
	"github.com/ethscanner/ethereum-log-scanner/core/mq"
	"github.com/ethscanner/ethereum-log-scanner/core/utils"
	"github.com/gogf/gf/v2/frame/g"
)

type sRmqLogStorage struct {
	logCache *cache.Cache
}

func NewgRmqLogStorage() *sRmqLogStorage {
	return &sRmqLogStorage{
		logCache: cache.NewCacheByName("sRmqLogStorage", 20*10000, nil), //预计可以缓存、1w+条记录
	}
}

func (s *sRmqLogStorage) SaveLogs(ctx context.Context, name string, logs []types.Log) error {
	for _, v := range logs {
		key := utils.FromatEventIdKey(name, v.BlockNumber, v.Index)
		if _, ok := s.logCache.Get(key); ok {
			g.Log().Infof(ctx, "key %v已经存在", key)
			continue
		}
		json, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if ok, err := mq.SendDefaultMessage(ctx, name, key, json); err != nil {
			return err
		} else if !ok {
			return errors.New("发送消息失败")
		} else {
			s.AddLogsToCache(ctx, name, v)
		}
	}
	return nil
}

func (s *sRmqLogStorage) AddLogsToCache(ctx context.Context, name string, log types.Log) error {
	key := utils.FromatEventIdKey(name, uint64(log.BlockNumber), log.Index)
	s.logCache.Add(key, cache.ByteView{})
	return nil
}
