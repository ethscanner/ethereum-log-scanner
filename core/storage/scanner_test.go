package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
)

const RPC_URL = "https://rpc.ankr.com/bsc"

func TestGetScanBlockNumbers(t *testing.T) {
	ctx := context.Background()
	store := NewGormScannerStorage()
	address := common.HexToAddress("")
	s := scanner.NewScanner("test", []common.Address{address}, store, nil, nil)
	var BlockNumber uint64 = 10000
	from, to, err := s.GetScanBlockNumbersByStroage(ctx, BlockNumber)
	if err != nil {
		t.Fatalf("error %e", err)
	}
	if to != BlockNumber-s.DelayBlocks {
		t.Fatalf("to number error")
	}
	if from != 0 {
		t.Fatalf("from number error")
	}
}

func TestFilterQuery(t *testing.T) {
	client, _ := ethclient.Dial(RPC_URL)
	ctx := context.Background()
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	address := common.HexToAddress("0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82")
	s := scanner.NewScanner("MarkTrasnferBlockNumber", []common.Address{address}, nil, nil, [][]common.Hash{{logTransferSigHash}})
	var start uint64 = 30635690
	var end uint64 = 30635691
	logs, err := s.FilterQuery(ctx, client, start, end)
	if err != nil {
		t.Fatalf("error %e", err)
	}
	fmt.Println(logs)

}

func TestConcurrentThreadQuery(t *testing.T) {
	client, _ := ethclient.Dial(RPC_URL)
	ctx := context.Background()
	address := common.HexToAddress("0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82")
	s := scanner.NewScanner("test", []common.Address{address}, nil, nil, nil)
	var start uint64 = 22921016
	var end uint64 = 22951016

	logs, scannedBlockNum, err := s.ConcurrentThreadQuery(ctx, client, start, end)
	if err != nil {
		t.Fatalf("error %e", err)
	}
	fmt.Println("logs length: ", len(logs))
	fmt.Println("scannedBlockNum: ", scannedBlockNum)
}

func TestSegmentationScan(t *testing.T) {
	client, _ := ethclient.Dial(RPC_URL)
	ctx := context.Background()
	logTransferSig := []byte("Transfer(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	fmt.Println(logTransferSigHash.Hex())
	logStorage := NewGormLogStorage()
	scanStorage := NewGormScannerStorage()
	address := common.HexToAddress("0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82")
	s := scanner.NewScanner("MarkTrasnferBlockNumber", []common.Address{address}, scanStorage, logStorage, [][]common.Hash{{logTransferSigHash}})
	var start uint64 = 26138478
	var end uint64 = 26139479

	logs, scannedBlockNum, err := s.SegmentationScan(ctx, client, start, end)
	if err != nil {
		t.Fatalf("error %o", err)
	}
	fmt.Println("logs length: ", len(logs))
	fmt.Println("scannedBlockNum: ", scannedBlockNum)

}

func TestScan(t *testing.T) {
	client, _ := ethclient.Dial(RPC_URL)
	ctx := context.Background()
	address := common.HexToAddress("0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82")
	_storage := NewGormLogStorage()
	scannerStorage := NewGormScannerStorage()
	s := scanner.NewScanner("MarkTrasnferBlockNumber", []common.Address{address}, scannerStorage, _storage, nil)
	_, err := s.ScanToStroage(ctx, client, 26139478)
	if err != nil {
		t.Fatalf("error %e", err)
	}
}
