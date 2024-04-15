package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestToml(t *testing.T) {
	// 首先初始化struct对象，用于保存配置文件信息
	var conf map[string]interface{}

	// 通过toml.DecodeFile将toml配置文件的内容，解析到struct对象
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		// handle error
	}
	path := "./config1.toml"
	exists, err := PathExists(path)
	if err != nil {
		panic(err)
	}
	var file *os.File
	if !exists {
		if file, err = os.Create(path); err != nil {
			panic(err)
		}
	}
	if file, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0644); err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	toml.NewEncoder(writer).Encode(conf)
	// 可以通过conf读取配置
	fmt.Println(conf)
}

func TestFileBase(t *testing.T) {
	files := "hello/config/config.tom"
	dir := filepath.Dir(files)
	fmt.Println(dir)
}
