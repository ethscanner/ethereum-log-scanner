// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
)

type (
	INftManagerHandle interface {
		// 启动检测
		ScanMintEvent(ctx context.Context) error
		// 处理事件
		HandleMint(ctx context.Context, events []scanner.Elog) error
	}
)

var (
	localNftManagerHandle INftManagerHandle
)

func NftManagerHandle() INftManagerHandle {
	if localNftManagerHandle == nil {
		panic("implement not found for interface INftManagerHandle, forgot register?")
	}
	return localNftManagerHandle
}

func RegisterNftManagerHandle(i INftManagerHandle) {
	localNftManagerHandle = i
}
