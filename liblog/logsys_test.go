package liblog

import (
	"fmt"
	"testing"
	"time"
)

func Testliblog(t *testing.T) {
	//全局日志
	InitLogConfig("./log/", "test", LOG_LEVEL_TRACE, 5, 10)

	LogInfo(0, "************** start ********************")
	RunLogInfo(0, "************** start ********************")

	if false {
		LogTrace(0, "a != b")
		LogDebug(0, "a != b")
		LogInfo(0, "a != b")
		LogWarn(0, "a != b")
		LogError(0, "a != b")
		LogFatal(0, "a != b")

		RunLogTrace(0, "a != b")
		RunLogDebug(0, "a != b")
		RunLogInfo(0, "a != b")
		RunLogWarn(0, "a != b")
		RunLogError(0, "a != b")
		RunLogFatal(0, "a != b")
	}

	//多协程打印日志
	if true {
		var nRuntineNum = 5
		chs := make([]chan int, nRuntineNum)
		for i := 0; i < nRuntineNum; i++ {
			chs[i] = make(chan int)
			go func(ch chan int, mark int) {
				fmt.Println("LoggerTest runtine:", mark)

				start := time.Now()
				for {
					LogTrace(0, mark)
					LogDebug(0, mark)
					LogInfo(0, mark)
					LogWarn(0, mark)
					LogError(0, mark)
					LogFatal(0, mark)

					RunLogTrace(0, mark)
					RunLogDebug(0, mark)
					RunLogInfo(0, mark)
					RunLogWarn(0, mark)
					RunLogError(0, mark)
					RunLogFatal(0, mark)

					if time.Now().Sub(start).Seconds() > 10 {
						break
					}
				}
				ch <- mark
			}(chs[i], i)
		}
		for _, ch := range chs {
			<-ch
		}
	}

	//局部日志
	if false {
		logger := NewLogInstance()
		logger.InitLog("./log/", "MBO", LOG_LEVEL_DEBUG, 1, 10)

		logger.LogTrace(0, "a != b")
		logger.LogDebug(0, "a != b")
		logger.LogInfo(0, "a != b")
		logger.LogWarn(0, "a != b")
		logger.LogError(0, "a != b")
		logger.LogFatal(0, "a != b")

		logger.RunLogTrace(0, "a != b")
		logger.RunLogDebug(0, "a != b")
		logger.RunLogInfo(0, "a != b")
		logger.RunLogWarn(0, "a != b")
		logger.RunLogError(0, "a != b")
		logger.RunLogFatal(0, "a != b")
	}
}
