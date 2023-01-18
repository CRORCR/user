package contract

import (
	"fmt"
	"math"
	"time"

	"context"

	"github.com/CRORCR/user/internal/model"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewLogger(logConfig model.LogConfig) {
	path := "./logs/%Y%m%d%H.log"
	writer, _ := retalog.New(
		path,
		retalog.WithMaxAge(time.Hour*24*time.Duration(logConfig.MaxDay)),
		retalog.WithRotationTime(time.Hour),
	)

	logrus.SetOutput(writer)

	// todo
	// 如果env是本地，控制台输出日志
	// 如果env是prod，只能输出info级别日志
	// 出现错误，比如panic的时候，需要打印堆栈信息 后面调研一下zap和glog，这两个包，好像是支持堆栈打印的

	//logrus.SetOutput(os.Stderr) // 控制台输出
	//logrus.SetOutput(ioutil.Discard) //控制台不输出

	switch logConfig.Level {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	//logrus.SetReportCaller(true) // 行号是否输出
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	t := time.Now()
	logrus.Infof(fmt.Sprintf("before handling request method [%s],request:[%v]", info.FullMethod, req))

	resp, err := handler(ctx, req)

	logrus.Infof(fmt.Sprintf("before handling request method [%s],response:[%v],time:[%d ms]", info.FullMethod, resp, int(math.Ceil(float64(time.Since(t).Nanoseconds()/1000000)))))
	return resp, err
}
