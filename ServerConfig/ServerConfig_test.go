package ServerConfig

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
	"time"
)

//示例xml文件:运行时将此文件内容保存
/*
<?xml version ="1.0" encoding="utf-8"?>
<config ver="1.0.2.1">
  <logpara echo="日志参数配置">
    <log_level dynamic="yes" echo="日志级别:0:TRACE;1:DEBUG;2:INFO;3:WARN;4:ERROR;5:FATAL;6:OFF">1</log_level>
    <log_file_size dynamic="no" echo="日志文件大小(MB)">10</log_file_size>
    <log_file_num dynamic="no" echo="日志文件个数">50</log_file_num>
    <log_path dynamic="no" echo="日志文件路径">log</log_path>
  </logpara>
  <db_info echo="msyql数据库配置">
    <db_server dynamic="no" echo="数据库地址">192.168.1.133</db_server>
    <db_port dynamic="no" echo="数据库端口">3306</db_port>
    <db_name dynamic="no" echo="数据库名称">mbo_server_db</db_name>
    <db_user dynamic="no" echo="数据库用户名">root</db_user>
    <db_passwd dynamic="no" echo="数据库密码">Suitang@20170601</db_passwd>
    <db_charset dynamic="no" echo="数据库客户端字符集">utf8mb4</db_charset>
    <db_max_conn dynamic="no" echo="数据库客户端最大连接数">8</db_max_conn>
  </db_info>
</config>

*/

type SrvConfig struct {
	BaseConfig
	Config ConfigInfo
}

type ConfigInfo struct {
	XMLName  xml.Name  `xml:"config"`
	LogParam LogConfig `xml:"logpara"`
	DbParam  DbConfig  `xml:"db_info"`
}

type LogConfig struct {
	XMLName     xml.Name `xml:"logpara"`
	LogLevel    int      `xml:"log_level"`
	LogFileSize int      `xml:"log_file_size"`
	LogFileNum  int      `xml:"log_file_num"`
	LogPath     string   `xml:"log_path"`
}

type DbConfig struct {
	XMLName   xml.Name `xml:"db_info"`
	DbServer  string   `xml:"db_server"`
	DbPort    int      `xml:"db_port"`
	DbName    string   `xml:"db_name"`
	DbUser    string   `xml:"db_user"`
	DbPasswd  string   `xml:"db_passwd"`
	DbCharset string   `xml:"db_charset"`
	DbMaxConn string   `xml:"db_max_conn"`
}

//读取配置信息
func (pConfig *SrvConfig) LoadItems(content []byte) error {
	fmt.Printf("SrvConfig::LoadItems  size:%d\n", len(content))

	err := xml.Unmarshal(content, &pConfig.Config)
	if err != nil {
		fmt.Println("Parse file fail. error:", err)
		return err
	}
	fmt.Printf("Parse success. Config:\n%v\n", pConfig.Config)
	return nil
}

/****************************** 执行测试 ********************************/

//时刻检测配置文件变化
func DetectConfig(ch chan int, sConfig IConfig) {
	start := time.Now()
	for time.Now().Sub(start).Seconds() < 10 {
		//定时监测配置文件
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

		time.Sleep(5 * time.Second)
	}
	fmt.Println("Quit runtine")
	ch <- 1
}

func TestMain(t *testing.T) {
	fmt.Printf("main pid:%d ppid:%d\n", os.Getpid(), os.Getppid())
	sConfig := new(SrvConfig)
	//var ch chan int
	ch := make(chan int)
	sConfig.Init("D:\\tp\\execute\\GO_PROJ\\src\\MainTest\\config\\config.xml")
	go DetectConfig(ch, sConfig)
	<-ch
	fmt.Println("Quit.")
}
