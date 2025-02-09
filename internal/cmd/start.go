package cmd

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethscanner/ethereum-log-scanner/core/mq"
	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	"github.com/ethscanner/ethereum-log-scanner/core/storage"
	"github.com/ethscanner/ethereum-log-scanner/internal/config"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
)

var mode string
var rpc string
var client *ethclient.Client
var lastBlock uint64

func Start(ctx context.Context) {
	config, err := config.ParseContractConfig(ctx)
	if err != nil {
		panic(err)
	}
	mode = config.Mode
	rpc = config.Rpc
	contracts := config.Contracts
	client, _ = ethclient.Dial(rpc)
	//解析配置文件
	g.Log().Infof(ctx, "配置信息: %v", contracts)
	if _, err := startQueryLastBockNum(ctx, client); err != nil {
		panic(err)
	}

	if mode == string(scanner.GORM_MODE) {
		startGorm(ctx, client, contracts)
	} else if mode == string(scanner.RMQ_MODE) {
		startRMQ(ctx, client, contracts)
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

func startGorm(ctx context.Context, client *ethclient.Client, listContractConfig []*config.ContractConfig) {
	logStorage := storage.NewGormLogStorage()
	scanStorage := storage.NewTomlScannerStorage()
	g.Log().Info(ctx, "startGorm")
	for _, v := range listContractConfig {
		_, err := startSingleContract(ctx, client, v, logStorage, scanStorage)
		if err != nil {
			panic(err)
		} else {
			g.Log().Infof(ctx, "%s 扫描启动成功", v.Name)
		}
		if !v.Check {
			continue
		}
		_, err = startCheckerSingleContract(ctx, client, v, logStorage)
		if err != nil {
			panic(err)
		} else {
			g.Log().Infof(ctx, "%s checker 启动成功", v.Name)
		}
	}
}

func startRMQ(ctx context.Context, client *ethclient.Client, listContractConfig []*config.ContractConfig) {
	mq.InitMQ(ctx)
	logStorage := storage.NewgRmqLogStorage()
	scanStorage := storage.NewTomlScannerStorage()
	g.Log().Info(ctx, "startGorm")
	for _, v := range listContractConfig {
		_, err := startSingleContract(ctx, client, v, logStorage, scanStorage)
		if err != nil {
			panic(err)
		} else {
			g.Log().Infof(ctx, "%s 扫描启动成功", v.Name)
		}
	}
}

func startSingleContract(ctx context.Context, client *ethclient.Client, config *config.ContractConfig, logStorage scanner.LogStorage, scanStorage scanner.ScannerStorage) (*gcron.Entry, error) {
	scan := scanner.NewScanner(config.Name, config.AddressObj, scanStorage, logStorage, config.TopicsObj)
	return gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		g.Log().Infof(ctx, "扫描%v开始*********************", config.Name)
		_, err := scan.ScanToStroage(ctx, client, lastBlock)
		if err != nil {
			g.Log().Infof(ctx, "扫描%v错误 %v", config.Name, err)
		}
		g.Log().Infof(ctx, "扫描%v结束---------------------", config.Name)
	})
}

func startCheckerSingleContract(ctx context.Context, client *ethclient.Client, config *config.ContractConfig, logStorage scanner.DbLogStorage) (*gcron.Entry, error) {
	checker := scanner.NewChecker(config.Name, logStorage)
	return gcron.AddSingleton(ctx, "*/10 * * * * *", func(ctx context.Context) {
		g.Log().Infof(ctx, "check %v开始*********************", config.Name)
		_, err := checker.CheckAllStroage(ctx, client, lastBlock)
		if err != nil {
			g.Log().Infof(ctx, "check %v错误 %v", config.Name, err)
		}
		g.Log().Infof(ctx, "check%v结束---------------------", config.Name)

	})
}
