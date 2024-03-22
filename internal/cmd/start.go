package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethscanner/ethereum-log-scanner/core/mq"
	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	"github.com/ethscanner/ethereum-log-scanner/core/storage"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
)

type ContractConfig struct {
	Name       string     `json:"name" gencodec:"required"`
	Check      bool       `json:"check" gencodec:"required"`
	Address    []string   `json:"address" gencodec:"required"`
	Topics     [][]string `json:"topics" gencodec:"required"`
	AddressObj []common.Address
	TopicsObj  [][]common.Hash
}

var listContractConfig []*ContractConfig
var mode string
var rpc string
var client *ethclient.Client
var lastBlock uint64
var logStorage scanner.LogStorage
var scanStorage scanner.ScannerStorage

func Start(ctx context.Context) {

	mode = g.Cfg().MustGet(ctx, "e-scanner.mode").String()
	rpc = g.Cfg().MustGet(ctx, "e-scanner.rpc").String()
	client, _ = ethclient.Dial(rpc)

	v, err := g.Cfg().Get(ctx, "e-scanner.contracts")
	if err != nil {
		panic(err)
	}
	err = v.Structs(&listContractConfig)
	if err != nil {
		panic(err)
	}
	//解析配置文件
	for _, v := range listContractConfig {
		topics := v.Topics
		v.TopicsObj = make([][]common.Hash, 0, len(topics))
		for index, topicLv1 := range topics {
			topicLv2Arr := make([]common.Hash, 0, len(topicLv1))
			for _, topicLv2 := range topicLv1 {
				if index == 0 {
					logTransferSig := []byte(topicLv2)
					hash := crypto.Keccak256Hash(logTransferSig)
					topicLv2Arr = append(topicLv2Arr, hash)
				} else {
					hash := common.HexToHash(topicLv2)
					topicLv2Arr = append(topicLv2Arr, hash)
				}
			}
			v.TopicsObj = append(v.TopicsObj, topicLv2Arr)
		}
		addresses := v.Address
		v.AddressObj = make([]common.Address, 0, len(addresses))
		for _, addrStr := range addresses {
			v.AddressObj = append(v.AddressObj, common.HexToAddress(addrStr))
		}
	}
	js, _ := json.Marshal(listContractConfig)
	g.Log().Infof(ctx, "配置信息: %v", string(js))
	startQueryLastBockNum(ctx, client)
	if mode == string(scanner.GORM_MODE) {
		startGorm(ctx, client, listContractConfig)
	} else if mode == string(scanner.RMQ_MODE) {
		startRMQ(ctx, client, listContractConfig)
	}
	g.Log().Info(ctx, "启动成功")
	select {}
}

func startQueryLastBockNum(ctx context.Context, client *ethclient.Client) (*gcron.Entry, error) {
	return gcron.AddSingleton(ctx, "*/5 * * * * *", func(ctx context.Context) {
		blockNumber, err := client.BlockNumber(ctx)
		if err != nil {
			g.Log().Infof(ctx, "最新区块错误 %v", err)
		} else {
			lastBlock = blockNumber
			g.Log().Infof(ctx, "最新区块: %d", blockNumber)
		}
	})
}

func startGorm(ctx context.Context, client *ethclient.Client, listContractConfig []*ContractConfig) {
	_logStorage := storage.NewGormLogStorage()
	logStorage = _logStorage
	scanStorage = storage.NewTomlScannerStorage()
	g.Log().Info(ctx, "startGorm")
	for _, v := range listContractConfig {
		_, err := startSingleContract(ctx, client, v)
		if err != nil {
			panic(err)
		} else {
			g.Log().Infof(ctx, "%s 扫描启动成功", v.Name)
		}
		if !v.Check {
			continue
		}
		_, err = startCheckerSingleContract(ctx, client, v, _logStorage)
		if err != nil {
			panic(err)
		} else {
			g.Log().Infof(ctx, "%s checker 启动成功", v.Name)
		}
	}
}

func startRMQ(ctx context.Context, client *ethclient.Client, listContractConfig []*ContractConfig) {
	mq.InitMQ(ctx)
	logStorage = storage.NewgRmqLogStorage()
	scanStorage = storage.NewTomlScannerStorage()
	g.Log().Info(ctx, "startGorm")
	for _, v := range listContractConfig {
		_, err := startSingleContract(ctx, client, v)
		if err != nil {
			panic(err)
		} else {
			g.Log().Infof(ctx, "%s 扫描启动成功", v.Name)
		}
	}
}

func startSingleContract(ctx context.Context, client *ethclient.Client, config *ContractConfig) (*gcron.Entry, error) {
	scan := scanner.NewScanner(config.Name, config.AddressObj, scanStorage, logStorage, config.TopicsObj)
	return gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		g.Log().Infof(ctx, "扫描%v开始*********************", config.Name)
		_, err := scan.ScanToStroage(ctx, client, lastBlock)
		if err != nil {
			fmt.Printf("err: %v \n", err)
		}
		g.Log().Infof(ctx, "扫描%v结束---------------------", config.Name)
	})
}

func startCheckerSingleContract(ctx context.Context, client *ethclient.Client, config *ContractConfig, _logStorage scanner.DbLogStorage) (*gcron.Entry, error) {
	checker := scanner.NewChecker(config.Name, _logStorage)
	return gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		g.Log().Infof(ctx, "check %v开始*********************", config.Name)
		checker.CheckAllStroage(ctx, client, lastBlock)
		g.Log().Infof(ctx, "check%v结束---------------------", config.Name)
	})
}
