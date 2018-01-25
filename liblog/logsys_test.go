package liblog

import (
	"fmt"
	"testing"
	"time"
)

func Testliblog(t *testing.T) {
	//全局日志
	InitLogConfig("./log/", "test", LOG_LEVEL_TRACE, 5, 10)

	LogInfo("************** start ********************")
	RunLogInfo("************** start ********************")

	if false {
		LogTrace("a != b")
		LogDebug("a != b")
		LogInfo("a != b")
		LogWarn("a != b")
		LogError("a != b")
		LogFatal("a != b")

		RunLogTrace("a != b")
		RunLogDebug("a != b")
		RunLogInfo("a != b")
		RunLogWarn("a != b")
		RunLogError("a != b")
		RunLogFatal("a != b")
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
					LogTrace(mark)
					LogDebug(mark)
					LogInfo(mark)
					LogWarn(mark)
					LogError(mark)
					LogFatal(mark)

					RunLogTrace(mark)
					RunLogDebug(mark)
					RunLogInfo(mark)
					RunLogWarn(mark)
					RunLogError(mark)
					RunLogFatal(mark)

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

		logger.LogTrace("a != b")
		logger.LogDebug("a != b")
		logger.LogInfo("a != b")
		logger.LogWarn("a != b")
		logger.LogError("a != b")
		logger.LogFatal("a != b")

		logger.RunLogTrace("a != b")
		logger.RunLogDebug("a != b")
		logger.RunLogInfo("a != b")
		logger.RunLogWarn("a != b")
		logger.RunLogError("a != b")
		logger.RunLogFatal("a != b")
	}
}
