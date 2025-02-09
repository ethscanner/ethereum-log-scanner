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

type Config = map[string]uint64

const CONFIG_DIR = "config"

type gTomlScannerStorage struct {
	lock sync.Mutex
}

func NewTomlScannerStorage() *gTomlScannerStorage {
	if exists, err := utils.PathExists(CONFIG_DIR); err != nil {
		panic(err)
	} else if !exists {
		if err := os.Mkdir(CONFIG_DIR, os.ModePerm); err != nil {
			panic(err)
		}
	}
	return &gTomlScannerStorage{}
}

// 创建配置文件
func (s *gTomlScannerStorage) CreateConfigFile(ctx context.Context, key string) (filePaht string, err error) {
	path := filepath.Join(CONFIG_DIR, key+".conf")
	if exists, err := utils.PathExists(path); err != nil {
		return path, err
	} else if !exists {
		if _, err := os.Create(path); err != nil {
			return path, err
		}
	}
	return path, nil
}

func (s *gTomlScannerStorage) UpdateUint64(ctx context.Context, key string, val uint64) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	filePath, err := s.CreateConfigFile(ctx, key)
	if err != nil {
		return err
	}
	var file *os.File
	defer file.Close()
	if file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644); err != nil {
		return err
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	config := Config{}
	config[key] = val
	return toml.NewEncoder(writer).Encode(config)
}

func (s *gTomlScannerStorage) InsertUint64(ctx context.Context, key string, val uint64) error {
	return s.UpdateUint64(ctx, key, val)
}

func (s *gTomlScannerStorage) GetUint64ByKey(ctx context.Context, key string) (val uint64, err error) {
	filePath, err := s.CreateConfigFile(ctx, key)
	if err != nil {
		return 0, err
	}
	config := Config{}
	// 通过toml.DecodeFile将toml配置文件的内容，解析到struct对象
	if _, err = toml.DecodeFile(filePath, &config); err != nil {
		return 0, err
	}
	return config[key], nil
}
