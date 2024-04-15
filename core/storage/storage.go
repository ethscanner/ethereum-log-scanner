package storage

import (
	"sync"

	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	"github.com/ethscanner/ethereum-log-scanner/core/storage/log"
	"github.com/ethscanner/ethereum-log-scanner/core/storage/position"
)

var tomlScannerInstance scanner.ScannerStorage
var lock sync.Mutex

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
	lock.Lock()
	defer lock.Unlock()
	if tomlScannerInstance == nil {
		tomlScannerInstance = position.NewTomlScannerStorage()
	}
	return tomlScannerInstance
}
