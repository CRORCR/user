package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Timeout(d time.Duration) gin.HandlerFunc {
	// 使用函数回调
	return func(c *gin.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		// 执行业务逻辑前预操作：初始化超时context
		durationCtx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 使用next执行具体的业务逻辑
			c.Next()

			finish <- struct{}{}
		}()
		// 执行业务逻辑后操作
		select {
		case p := <-panicChan:
			c.JSON(http.StatusInternalServerError, "Server time out")
			// 服务异常走到这里，无法打印行号，所以各个业务自己处理自己的panic
			logrus.Errorf("服务panic:%v", p)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.JSON(http.StatusInternalServerError, "Server time out")
			log.Println("服务超时")
		}
	}
}
