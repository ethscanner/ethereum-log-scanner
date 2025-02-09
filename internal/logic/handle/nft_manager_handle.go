package handle

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	"github.com/ethscanner/ethereum-log-scanner/core/storage"
	"github.com/ethscanner/ethereum-log-scanner/internal/dao"
	"github.com/ethscanner/ethereum-log-scanner/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

var nftManagerMintSig = []byte("Mint(address,address,uint256,uint256,uint256)")
var nftManagerMintSigHash = crypto.Keccak256Hash(nftManagerMintSig)

type sNftManagerHandle struct {
	logStorage scanner.DbLogStorage
	address    string
}

func init() {
	logStorage := storage.NewGormLogStorage()
	nftManagerAddr, err := g.Cfg().Get(gctx.GetInitCtx(), "contract.address.nftManager")
	if err != nil {
		panic(err)
	}
	service.RegisterNftManagerHandle(NewNftManagerHandle(logStorage, nftManagerAddr.String()))
}

func NewNftManagerHandle(logStorage scanner.DbLogStorage, address string) *sNftManagerHandle {
	return &sNftManagerHandle{
		logStorage: logStorage,
		address:    address,
	}
}

// 启动检测
func (s *sNftManagerHandle) ScanMintEvent(ctx context.Context) error {
	lastTime := gtime.Now().Add(-time.Hour * 24 * 30)
	var State int = int(scanner.LOG_STATE_PENDING)
	eventHash := nftManagerMintSigHash.String()
	query := scanner.LogQuery{
		ContractAddress: s.address,
		EventHash:       eventHash,
		State:           &State,
		CreatedAtGt:     lastTime,
		Limit:           1000,
	}
	if logs, err := s.logStorage.QueryLogs(ctx, query); err != nil {
		return err
	} else {
		if err := s.HandleMint(ctx, logs); err != nil {
			return err
		}

	}
	return nil

}

// 处理事件
func (s *sNftManagerHandle) HandleMint(ctx context.Context, events []scanner.Elog) error {
	var nftMintIds []uint64
	var logIds []uint64
	for _, log := range events {
		//hashs := utils.BytesToHashArray(log.Data)
		id := log.Log.Topics[3].Big().Uint64()
		nftMintIds = append(nftMintIds, id)
		logIds = append(logIds, log.Id)
		//fmt.Println(v.TxHash)
	}
	if len(nftMintIds) > 0 {
		if _, err := dao.DesertNftMintRecords.Ctx(ctx).
			Where("status = ? AND id IN (?)", scanner.LOG_STATE_PENDING, nftMintIds).Update(g.Map{
			"status": scanner.LOG_STATE_PROCESSED,
		}); err != nil {
			return err
		}

	}
	if len(logIds) > 0 {
		if err := s.logStorage.MarkAsProcessed(ctx, "nft_manager_mint", logIds); err != nil {
			return err
		}
	}
	return nil

}
