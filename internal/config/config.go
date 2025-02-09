package config

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogf/gf/v2/frame/g"
)

type ContractConfig struct {
	Name       string     `json:"name" gencodec:"required"`
	Check      bool       `json:"check" gencodec:"required"`
	Address    []string   `json:"address" gencodec:"required"`
	Topics     [][]string `json:"topics" gencodec:"required"`
	AddressObj []common.Address
	TopicsObj  [][]common.Hash
}

type ScannerConfig struct {
	Contracts []*ContractConfig `json:"contracts" gencodec:"required"`
	Mode      string            `json:"mode" gencodec:"required"`
	Rpc       string            `json:"rpc" gencodec:"required"`
}

func ParseContractConfig(ctx context.Context) (scannerConfig *ScannerConfig, err error) {
	scannerConfig = &ScannerConfig{}
	scannerConfig.Mode = g.Cfg().MustGet(ctx, "e-scanner.mode").String()
	scannerConfig.Rpc = g.Cfg().MustGet(ctx, "e-scanner.rpc").String()
	v, err := g.Cfg().Get(ctx, "e-scanner.contracts")
	if err != nil {
		panic(err)
	}
	err = v.Structs(&scannerConfig.Contracts)
	if err != nil {
		panic(err)
	}
	//解析配置文件
	for _, v := range scannerConfig.Contracts {
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

	return scannerConfig, nil
}
