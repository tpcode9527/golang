package liblog

/**/
import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

/* 日志类型 */
const (
	LOG_TYPE_SYS = 0 //系统日志
	LOG_TYPE_RUN = 1 //运行日志
)

/*日志等级*/
const (
	LOG_LEVEL_TRACE = 0
	LOG_LEVEL_DEBUG = 1
	LOG_LEVEL_INFO  = 2
	LOG_LEVEL_WARN  = 3
	LOG_LEVEL_ERROR = 4
	LOG_LEVEL_FATAL = 5
)

/*基本日志配置参数*/
type LogParam struct {
	Path     string //文件路径;
	Prefix   string //日志前缀用于区分模块
	LogLevel int    //日志等级
	LogSize  int64  //单个日志文件大小
	LogNum   int    //日志文件保留个数
}

/*基本日志文件信息*/
type LogFileInfo struct {
	path string
	file *os.File
	mtx  *sync.RWMutex
}

/*日志实例参数*/
type LogInstance struct {
	LogParm  LogParam //基本日志参数
	FileInfo map[int]*LogFileInfo
}

/*全局日志实例*/
var inst *LogInstance

/*日志等级等级*/
var mapLogLevel map[int]string

/*日志文件名前缀*/
var mapLogFilePrefix map[int]string

/*初始化创建默认的日志参数*/
func init() {
	//日志等级提示
	mapLogLevel = make(map[int]string)
	mapLogLevel[LOG_LEVEL_TRACE] = "TRACE"
	mapLogLevel[LOG_LEVEL_DEBUG] = "DEBUG"
	mapLogLevel[LOG_LEVEL_INFO] = "INFO"
	mapLogLevel[LOG_LEVEL_WARN] = "WARN"
	mapLogLevel[LOG_LEVEL_ERROR] = "ERROR"
	mapLogLevel[LOG_LEVEL_FATAL] = "FATAL"

	//日志类型
	mapLogFilePrefix = make(map[int]string)
	mapLogFilePrefix[LOG_TYPE_SYS] = "sys"
	mapLogFilePrefix[LOG_TYPE_RUN] = "run"

	//全局日志实例
	inst = NewLogInstance()
}

/*检测文件是否存在*/
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*获取文件大小*/
func GetFileSize(file string) int64 {
	sFileInfo, err := os.Stat(file)
	if nil != err {
		fmt.Println("Load config fail. error:", err)
		return -1
	}
	return sFileInfo.Size()
}

/*获取完整文件路径*/
func getLogFile(path string, logType int, prefix string, seq int) string {
	var logFile string = path + mapLogFilePrefix[logType] + prefix + ".log"

	if 0 != seq {
		logFile += "." + strconv.Itoa(seq)
	}

	return logFile
}

/*不安全的打开方式*/
func (this *LogFileInfo) openFile_unsafe() error {

	var err error = nil
	//打开系统日志文件
	if bResult, _ := PathExists(this.path); bResult {
		this.file, err = os.OpenFile(this.path, os.O_RDWR|os.O_APPEND, 0666)
		if nil != err {
			fmt.Println("Open file fail. error:", err)
			return err
		}
	} else {
		this.file, err = os.OpenFile(this.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if nil != err {
			fmt.Println("Create file fail. error:", err)
			return err
		}
	}
	return nil
}

/*打开文件*/
func (this *LogFileInfo) openFile() error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	return this.openFile_unsafe()
}

/*判定文件是io否有效*/
func (this *LogFileInfo) isFileValid_unsafe() bool {
	return (nil != this.file)
}

/*安全判定文件是io否有效*/
func (this *LogFileInfo) isFileValid() bool {
	this.mtx.RLock()
	defer this.mtx.RUnlock()

	return this.isFileValid_unsafe()
}

/*保存内容*/
func (this *LogFileInfo) writeFile(a ...interface{}) error {
	this.mtx.RLock()
	defer this.mtx.RUnlock()

	if !this.isFileValid_unsafe() {
		fmt.Println("File is not openning.")
		return errors.New("File is not openning.")
	}

	_, err := fmt.Fprintln(this.file, a)
	/*
			text := fmt.Sprintf("%s\n", a)
			_, err := this.file.WriteString(text)
			if nil != err {
			fmt.Println("WriteString fail. error: ", err)
		} else {
			this.file.Sync()
		}
	*/

	return err
}

/*关闭文件*/
func (this *LogFileInfo) closeFile_unsafe() error {
	if this.isFileValid_unsafe() {
		this.file.Close()
		this.file = nil
	}

	return nil
}

/*安全关闭文件*/
func (this *LogFileInfo) closeFile() error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	return this.closeFile_unsafe()
}

/*重命名文件*/
func (this *LogFileInfo) renameFile(dstFile string) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	err := os.Rename(this.path, dstFile)
	if nil != err {
		fmt.Println("Rename ", this.path, " to ", dstFile, " fail. error: ", err)
		return err
	}

	return nil
}

/*滚动修改已经存在的日志文件*/
func (this *LogFileInfo) rollFile(path string, logType int, prefix string, fileNum int) error {
	//关闭文件
	this.closeFile_unsafe()

	//循环修改已经生成的日志文件
	var srcFile string
	var dstFile string
	for i := fileNum - 1; i >= 1; i-- {
		srcFile = getLogFile(path, logType, prefix, i)
		dstFile = getLogFile(path, logType, prefix, i+1)
		err1 := os.Rename(srcFile, dstFile)
		if nil != err1 {
			//fmt.Println("Rename fail. error: ", err1)
		}
	}

	//修改当前日志文件
	dstFile = getLogFile(path, logType, prefix, 1)
	srcFile = getLogFile(path, logType, prefix, 0)
	err2 := os.Rename(srcFile, dstFile)
	if nil != err2 {
		fmt.Println("Rename fail. error: ", err2)
	}

	//重新打开文件
	err := this.openFile_unsafe()
	if nil != err {
		return err
	}

	return nil
}

/*记录日志*/
func (this *LogFileInfo) writeLog(path string, logType int, prefix string, fileNum int, fileSize int64, a ...interface{}) error {
	err := this.writeFile(a)
	if nil != err {
		return err
	}

	this.mtx.Lock()
	defer this.mtx.Unlock()

	if GetFileSize(this.path) >= fileSize {
		err = this.rollFile(path, logType, prefix, fileNum)
	}

	return err
}

/*打开日志文件*/
func (this *LogInstance) openLogFile() error {
	var err error = nil

	//打开系统日志文件
	err = this.FileInfo[LOG_TYPE_SYS].openFile()
	if nil != err {
		return err
	}

	//打开运行日志文件
	err = this.FileInfo[LOG_TYPE_RUN].openFile()
	if nil != err {
		return err
	}

	return nil
}

/*新建日志实例*/
func NewLogInstance() *LogInstance {
	return &LogInstance{LogParm: LogParam{LogLevel: LOG_LEVEL_TRACE, LogSize: 50, LogNum: 5}, FileInfo: make(map[int]*LogFileInfo)}
}

/*获取完整文件路径*/
func (this *LogInstance) getFilePath(logType int, seq int) string {
	return getLogFile(this.LogParm.Path, logType, this.LogParm.Prefix, seq)
}

/*初始化日志文件配置*/
func (this *LogInstance) InitLog(path string, prefix string, logLevel int, logSize int64, logNum int) error {
	//纠正日志级别
	if logLevel > LOG_LEVEL_FATAL {
		logLevel = LOG_LEVEL_FATAL
	} else if logLevel < LOG_LEVEL_TRACE {
		logLevel = LOG_LEVEL_TRACE
	}

	//日志文件设置
	if logSize < 1 || logNum < 1 {
		fmt.Println("Init log fail.")
		return errors.New("Init log fail.")
	}

	//设置日志配置
	if 0 != len(path) {
		if exist, _ := PathExists(path); !exist {
			os.MkdirAll(path, 0777)
		}

		if len(path)-1 == strings.LastIndex(path, "/") {
			this.LogParm.Path = path
		} else {
			this.LogParm.Path = path + "/"
		}
	}
	this.LogParm.Prefix = prefix
	this.LogParm.LogLevel = logLevel
	this.LogParm.LogSize = logSize * 1024 * 1024 //以M为单位
	this.LogParm.LogNum = logNum

	//日志文件
	this.FileInfo = make(map[int]*LogFileInfo)
	this.FileInfo[LOG_TYPE_SYS] = &LogFileInfo{file: nil, mtx: new(sync.RWMutex)}
	this.FileInfo[LOG_TYPE_SYS].path = this.getFilePath(LOG_TYPE_SYS, 0)

	this.FileInfo[LOG_TYPE_RUN] = &LogFileInfo{file: nil, mtx: new(sync.RWMutex)}
	this.FileInfo[LOG_TYPE_RUN].path = this.getFilePath(LOG_TYPE_RUN, 0)

	//打开系统日志文件
	return this.openLogFile()
}

/*向日志文件中写入日志*/
func (this *LogInstance) writeLogFile(logType int, a ...interface{}) error {
	return this.FileInfo[logType].writeLog(this.LogParm.Path, logType, this.LogParm.Prefix, this.LogParm.LogNum, this.LogParm.LogSize, a)
}

/*
* 保存日志信息
 */
func (this *LogInstance) SaveLog(logType int, logLevel int, a ...interface{}) error {
	//对于配置的日志等级以下的日志信息不予显示
	if logLevel < inst.LogParm.LogLevel {
		return nil
	}

	fileInfo, ok := this.FileInfo[logType]
	if !ok || !fileInfo.isFileValid() {
		fmt.Println("File is not open or config missing. or:", ok)
		_, err := fmt.Println(time.Now().Format("2006-01-02 15:04:05"), mapLogLevel[logLevel], a)
		return err
	}

	return this.writeLogFile(logType, time.Now().Format("2006-01-02 15:04:05"), mapLogLevel[logLevel], a)
}

/********************** 系统日志 **************************/
func (this *LogInstance) LogTrace(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_SYS, LOG_LEVEL_TRACE, a)
}

func (this *LogInstance) LogDebug(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_SYS, LOG_LEVEL_DEBUG, a)
}

func (this *LogInstance) LogInfo(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_SYS, LOG_LEVEL_INFO, a)
}

func (this *LogInstance) LogWarn(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_SYS, LOG_LEVEL_WARN, a)
}

func (this *LogInstance) LogError(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_SYS, LOG_LEVEL_ERROR, a)
}

func (this *LogInstance) LogFatal(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_SYS, LOG_LEVEL_FATAL, a)
}

/********************** 运行日志 **************************/
func (this *LogInstance) RunLogTrace(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_RUN, LOG_LEVEL_TRACE, a)
}

func (this *LogInstance) RunLogDebug(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_RUN, LOG_LEVEL_DEBUG, a)
}

func (this *LogInstance) RunLogInfo(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_RUN, LOG_LEVEL_INFO, a)
}

func (this *LogInstance) RunLogWarn(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_RUN, LOG_LEVEL_WARN, a)
}

func (this *LogInstance) RunLogError(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_RUN, LOG_LEVEL_ERROR, a)
}

func (this *LogInstance) RunLogFatal(a ...interface{}) error {
	return this.SaveLog(LOG_TYPE_RUN, LOG_LEVEL_FATAL, a)
}

/*
* 初始化日志配置
* path 日志文件保存路径
* prefix 日志文件前缀，用于不同模块区分，如传入的是MBO。那么日志文件就是sysMBO.log runMBO.log
* logLevel 要保存的日志等级
* logSize 每个日志文件大小
* logNum 保存的日志文件总数
 */
func InitLogConfig(path string, prefix string, logLevel int, logSize int64, logNum int) error {
	return inst.InitLog(path, prefix, logLevel, logSize, logNum)
}

/********************** 系统日志 **************************/
func LogTrace(a ...interface{}) error {
	return inst.LogTrace(a)
}

func LogDebug(a ...interface{}) error {
	return inst.LogDebug(a)
}

func LogInfo(a ...interface{}) error {
	return inst.LogInfo(a)
}

func LogWarn(a ...interface{}) error {
	return inst.LogWarn(a)
}

func LogError(a ...interface{}) error {
	return inst.LogError(a)
}

func LogFatal(a ...interface{}) error {
	return inst.LogFatal(a)
}

/********************** 运行日志 **************************/
func RunLogTrace(a ...interface{}) error {
	return inst.RunLogTrace(a)
}

func RunLogDebug(a ...interface{}) error {
	return inst.RunLogDebug(a)
}

func RunLogInfo(a ...interface{}) error {
	return inst.RunLogInfo(a)
}

func RunLogWarn(a ...interface{}) error {
	return inst.RunLogWarn(a)
}

func RunLogError(a ...interface{}) error {
	return inst.RunLogError(a)
}

func RunLogFatal(a ...interface{}) error {
	return inst.RunLogFatal(a)
}
