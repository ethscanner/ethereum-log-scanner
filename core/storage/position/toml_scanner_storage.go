package position

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/ethscanner/ethereum-log-scanner/core/utils"
)

const FILE_PATH = "config/scanner.toml"

type Config = map[string]uint64
type gTomlScannerStorage struct {
	lock sync.Mutex
}

func NewTomlScannerStorage() *gTomlScannerStorage {
	dir := filepath.Dir(FILE_PATH)
	if exists, err := utils.PathExists(dir); err != nil {
		panic(err)
	} else if !exists {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	if exists, err := utils.PathExists(FILE_PATH); err != nil {
		panic(err)
	} else if !exists {
		if _, err := os.Create(FILE_PATH); err != nil {
			panic(err)
		}
	}
	return &gTomlScannerStorage{}
}

func (s *gTomlScannerStorage) GetConfig(ctx context.Context) (conf Config, err error) {
	// 通过toml.DecodeFile将toml配置文件的内容，解析到struct对象
	if _, err = toml.DecodeFile(FILE_PATH, &conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (s *gTomlScannerStorage) UpdateUint64(ctx context.Context, key string, val uint64) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	config, err := s.GetConfig(ctx)
	if err != nil {
		return err
	}
	var file *os.File
	if file, err = os.OpenFile(FILE_PATH, os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644); err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	config[key] = val
	return toml.NewEncoder(writer).Encode(config)
}

func (s *gTomlScannerStorage) InsertUint64(ctx context.Context, key string, val uint64) error {
	return s.UpdateUint64(ctx, key, val)
}

func (s *gTomlScannerStorage) GetUint64ByKey(ctx context.Context, key string) (val uint64, err error) {

	config, err := s.GetConfig(ctx)
	if err != nil {
		return 0, err
	}
	return config[key], nil
}
