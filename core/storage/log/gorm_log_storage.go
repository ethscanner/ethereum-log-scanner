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
)

type gOrmLogStorage struct {
	logCache *cache.Cache
}

type HdContractEventObj struct {
	ContractName    string `json:"contractName"    ` // 合约名
	ContractAddress string `json:"contractAddress" ` // 合约地址
	TxHash          string `json:"txHash"          ` // 交易哈希
	EventHash       string `json:"eventHash"       ` // 事件名
	EventId         int64  `json:"eventId"         ` // 事件id
	BlockNumber     int64  `json:"blockNumber"     ` // 区块编号
	BlockHash       string `json:"blockHash"       ` // 交易哈希
	EventTopics     []byte `json:"eventTopics"     ` // 事件头
	EventData       []byte `json:"eventData"       ` // event数据
	State           int    `json:"state"           ` // 0:待处理 10:已处理

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
	saveLogs := make([]*HdContractEventObj, 0, len(logs))
	for _, v := range logs {
		txHash := v.TxHash.Hex()
		EventId := int64(v.Index)
		key := utils.FromatEventIdKey(name, v.BlockHash.Hex(), v.Index)
		if _, ok := s.logCache.Get(key); ok {
			g.Log().Infof(ctx, "key %v已经存在", key)
			continue
		}
		obj := HdContractEventObj{
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
		}
		saveLogs = append(saveLogs, &obj)
	}
	ret, err := sdao.Data(saveLogs).Batch(5000).InsertIgnore()
	if err != nil {
		return err
	} else {
		if rows, _ := ret.RowsAffected(); rows != int64(len(saveLogs)) {
			g.Log().Errorf(ctx, "插入数量:%d,插入成功数量:%d", len(saveLogs), rows)
		} else {
			return s.AddLogsToCache(ctx, name, saveLogs)
		}
	}
	return nil
}

func (s *gOrmLogStorage) AddLogsToCache(ctx context.Context, name string, logs []*HdContractEventObj) error {
	for _, v := range logs {
		key := utils.FromatEventIdKey(name, v.BlockHash, uint(v.EventId))
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
	if query.State != nil {
		sdao = sdao.Where("state=?", query.State)
	}
	if len(query.StateList) > 0 {
		sdao = sdao.Where("state in (?)", query.StateList)
	}
	if query.BlockNumber != nil {
		sdao = sdao.Where("block_number=?", query.BlockNumber)
	}
	if query.BlockNumberLt != nil {
		sdao = sdao.Where("block_number < ?", query.BlockNumberLt)
	}
	if query.BlockNumberGt != nil {
		sdao = sdao.Where("block_number > ?", query.BlockNumberGt)
	}
	if query.EventId != nil {
		sdao = sdao.Where("event_id=?", query.EventId)
	}
	if query.CreatedAtGt != nil {
		sdao = sdao.Where("created_at > ?", query.CreatedAtGt)
	}
	if query.CreatedAtLt != nil {
		sdao = sdao.Where("created_at < ?", query.CreatedAtLt)
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
	if query.CheckState != nil {
		sdao = sdao.Where("check_state=?", query.CheckState)
	}
	if len(query.CheckStateList) > 0 {
		sdao = sdao.Where("check_state in (?)", query.CheckStateList)
	}
	if query.IdGt != nil {
		sdao = sdao.Where("id > ?", query.IdGt)
	}
	if query.Limit == 0 || query.Limit > 10000 {
		query.Limit = 10000
	}
	dbList := []entity.HdContractEvent{}
	orderBy := "block_number,event_id"
	if query.OrderBy != "" {
		orderBy = query.OrderBy
	}
	if query.Desc {
		sdao = sdao.OrderDesc(orderBy)
	} else {
		sdao = sdao.OrderAsc(orderBy)
	}
	err = sdao.Limit(query.Limit).Scan(&dbList)
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
		ContractName: v.ContractName,
		Id:           v.Id,
		Log:          tlog,
		CheckState:   v.CheckState,
		CheckedBlock: v.CheckedBlock,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

func (s *gOrmLogStorage) UpdateBlockCheckState(ctx context.Context, blockHashs []string, blockNumber uint64, checkState int) error {
	sdao := dao.HdContractEvent.Ctx(ctx)
	_, err := sdao.Data("check_state = ?,checked_block=?", checkState, blockNumber).
		Where("block_hash in (?)", blockHashs).Update()
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
