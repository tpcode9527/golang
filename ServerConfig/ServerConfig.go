package ServerConfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//接口
type IConfig interface {
	Init(path string)
	IsModify() bool
	ReadContent() ([]byte, error)
	LoadItems(content []byte) error
}

//基本配置信息
type BaseConfig struct {
	FilePath string
	LastTime time.Time
}

//初始化配置文件路径
func (pConfig *BaseConfig) Init(path string) {
	pConfig.FilePath = path
}

//判定文件是否被修改
func (pConfig *BaseConfig) IsModify() bool {
	sFileInfo, err := os.Stat(pConfig.FilePath)
	if nil != err {
		fmt.Println("Load config fail. error:", err)
		return false
	}

	fileTime := sFileInfo.ModTime()
	if pConfig.LastTime != fileTime {
		pConfig.LastTime = fileTime
		return true
	}

	return false
}

//读取文件信息
func (pConfig *BaseConfig) ReadContent() ([]byte, error) {
	content, err := ioutil.ReadFile(pConfig.FilePath)
	if err != nil {
		fmt.Println("Read file fail. error:", err)
	}
	return content, err
}

//读取基本配置项目信息
func (pConfig *BaseConfig) LoadItems(content []byte) error {
	fmt.Println("BaseConfig::LoadItems")
	return nil
}
