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
	ParseConfig() error
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

/*
* 解析配置文件
 */
func (this *BaseConfig) ParseConfig() error {
	//读取文件数据
	content, err := this.ReadContent()
	if err != nil {
		fmt.Println("Read fail. error:", err)
	} else {

		//解析配置文件内容
		if err = this.LoadItems(content); nil != err {
			fmt.Println("Parse config fail. error:", err)
		}
	}
	return err
}

/*
//时刻检测配置文件变化示例
func DetectConfig(chStop chan int, chEnd chan int, sConfig IConfig) {
	for {
		select {
		case <-chStop:
			fmt.Println("Receive stop channel")
			break
		default:
		}

		//定时监测配置文件
		time.Sleep(50 * time.Millisecond)
		//fmt.Printf("Detect config file. runtine pid:%d ppid:%d\n", os.Getpid(), os.Getppid())

		//检测是否有文件变化
		if sConfig.IsModify() {

			//读取文件数据
			content, err := sConfig.ReadContent()
			if err != nil {
				fmt.Println("Read fail. error:", err)
			} else {

				//解析配置文件内容
				if err = sConfig.LoadItems(content); nil != err {
					fmt.Println("Parse config fail. error:", err)
				}
			}
		}
	}
	fmt.Println("Quit DetectConfig")
	chEnd <- 1
}
*/
