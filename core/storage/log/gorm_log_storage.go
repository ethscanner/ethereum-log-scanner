package log

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethscanner/ethereum-log-scanner/core/cache"
	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	"github.com/ethscanner/ethereum-log-scanner/core/utils"
	"github.com/ethscanner/ethereum-log-scanner/internal/dao"
	"github.com/ethscanner/ethereum-log-scanner/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type gOrmLogStorage struct {
	logCache *cache.Cache
}

func NewGormLogStorage() *gOrmLogStorage {
	return &gOrmLogStorage{
		logCache: cache.NewCacheByName("gOrmLogStorage", 20*10000, nil), //预计可以缓存、1w+条记录
	}
}

func (s *gOrmLogStorage) SaveLogs(ctx context.Context, name string, logs []types.Log) error {
	if logs == nil || len(logs) == 0 {
		return nil
	}
	sdao := dao.HdContractEvent.Ctx(ctx)
	batch := len(logs) / 5000
	saveLogs := make([]entity.HdContractEvent, 0, len(logs))
	for _, v := range logs {
		txHash := v.TxHash.Hex()
		EventId := int64(v.Index)
		key := utils.FromatEventIdKey(name, v.BlockNumber, v.Index)
		if _, ok := s.logCache.Get(key); ok {
			g.Log().Infof(ctx, "key %v已经存在", key)
			continue
		}
		obj := entity.HdContractEvent{
			ContractName:    name,
			ContractAddress: v.Address.Hex(),
			TxHash:          txHash,
			EventHash:       v.Topics[0].Hex(),
			EventId:         EventId,
			BlockNumber:     int64(v.BlockNumber),
			BlockHash:       v.BlockHash.Hex(),
			EventTopics:     utils.HashArrayToBytes(v.Topics),
			EventData:       v.Data,
			State:           0,
			CreatedAt:       gtime.Now(),
		}
		saveLogs = append(saveLogs, obj)
	}
	if _, err := sdao.Data(saveLogs).Batch(batch).InsertIgnore(); err != nil {
		return err
		// } else if rows, _ := ret.RowsAffected(); int(rows) != len(saveLogs) {
		// 	return errors.New("保存日志错误")
	} else {
		return s.AddLogsToCache(ctx, name, saveLogs)
	}

}

func (s *gOrmLogStorage) AddLogsToCache(ctx context.Context, name string, logs []entity.HdContractEvent) error {
	for _, v := range logs {
		key := utils.FromatEventIdKey(name, uint64(v.BlockNumber), uint(v.EventId))
		s.logCache.Add(key, cache.ByteView{})
	}
	return nil
}

func (s *gOrmLogStorage) MarkAsProcessed(ctx context.Context, name string, ids []uint64) error {
	sdao := dao.HdContractEvent.Ctx(ctx)

	_, err := sdao.Data("state=10").WhereIn("id", ids).Update()
	return err
}

func (s *gOrmLogStorage) QueryLogs(ctx context.Context, query scanner.LogQuery) (logs []scanner.Elog, err error) {
	sdao := dao.HdContractEvent.Ctx(ctx)
	if query.State >= 0 {
		sdao = sdao.Where("state=?", query.State)
	}
	if query.BlockNumber != 0 {
		sdao = sdao.Where("block_number=?", query.BlockNumber)
	}
	if query.EventId != 0 {
		sdao = sdao.Where("event_id=?", query.EventId)
	}
	if query.ContractName != "" {
		sdao = sdao.Where("contract_name=?", query.ContractName)
	}
	if query.EventHash != "" {
		sdao = sdao.Where("event_hash=?", query.EventHash)
	}
	if query.ContractAddress != "" {
		sdao = sdao.Where("contract_address=?", query.ContractAddress)
	}
	if query.TxHash != "" {
		sdao = sdao.Where("tx_hash=?", query.TxHash)
	}
	if query.CheckState >= 0 {
		sdao = sdao.Where("check_state=?", query.CheckState)
	}
	if query.IdGt > 0 {
		sdao = sdao.Where("id > ?", query.IdGt)
	}
	if query.Limit == 0 || query.Limit > 10000 {
		query.Limit = 10000
	}
	dbList := []entity.HdContractEvent{}
	err = sdao.OrderDesc("id").Limit(query.Limit).Scan(&dbList)
	if err != nil {
		return
	}
	logs = make([]scanner.Elog, len(dbList))
	for i, v := range dbList {
		logs[i] = s.entity2Elog(v)
	}
	return
}

func (s *gOrmLogStorage) entity2Elog(v entity.HdContractEvent) (elog scanner.Elog) {
	tlog := types.Log{
		Address:     common.HexToAddress(v.ContractAddress),
		Topics:      utils.BytesToHashArray(v.EventTopics),
		Data:        v.EventData,
		BlockNumber: uint64(v.BlockNumber),
		BlockHash:   common.HexToHash(v.BlockHash),
		TxHash:      common.HexToHash(v.TxHash),
		Index:       uint(v.EventId),
	}
	return scanner.Elog{
		Id:           v.Id,
		Log:          tlog,
		CheckState:   v.CheckState,
		CheckedBlock: v.CheckedBlock,
	}
}

func (s *gOrmLogStorage) UpdateBlockCheckState(ctx context.Context, log scanner.Elog) error {
	sdao := dao.HdContractEvent.Ctx(ctx)
	_, err := sdao.Data("check_state = ?,checked_block=?", log.CheckState, log.CheckedBlock).
		Where("id=?", log.Id).Update()
	if err != nil {
		return err
	}
	return nil
}

func (s *gOrmLogStorage) GetLogByElogId(ctx context.Context, txHash common.Hash, blockHash common.Hash, index uint) (log scanner.Elog, err error) {
	ret, err := dao.HdContractEvent.Ctx(ctx).
		Where("tx_hash=? and block_hash=? and event_id=?",
			txHash.Hex(), blockHash.Hex(), index).One()
	if err != nil || ret.IsEmpty() {
		return log, err
	} else {
		entityObj := entity.HdContractEvent{}
		if err := ret.Struct(&entityObj); err != nil {
			return log, err
		}
		return s.entity2Elog(entityObj), nil
	}
}
