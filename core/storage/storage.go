package storage

import (
	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	"github.com/ethscanner/ethereum-log-scanner/core/storage/log"
	"github.com/ethscanner/ethereum-log-scanner/core/storage/position"
)

func NewGormLogStorage() scanner.DbLogStorage {
	return log.NewGormLogStorage()
}

func NewgRmqLogStorage() scanner.LogStorage {
	return log.NewgRmqLogStorage()
}

func NewGormScannerStorage() scanner.ScannerStorage {
	return position.NewGormScannerStorage()
}

func NewTomlScannerStorage() scanner.ScannerStorage {
	return position.NewTomlScannerStorage()
}
