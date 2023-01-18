package main

import (
	"fmt"
	"log"
	"net"

	callPrice "github.com/CRORCR/duoo-common/proto/call_price"
	"github.com/CRORCR/user/app/server"
	"github.com/CRORCR/user/internal/config"
	"github.com/CRORCR/user/internal/contract"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// 加载配置
	config := config.InitConfig()
	contract.NewLogger(config.Conf.Log)

	addr := fmt.Sprintf(":%s", config.Conf.Rpc.Port)

	log.Printf("server listen address in %s \n", addr)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Conf.Rpc.Port))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = l.Close(); err != nil {
			logrus.Errorf("Failed to close %s %s: %v", "tcp", addr, err)
		}
	}()

	s := grpc.NewServer(grpc.UnaryInterceptor(contract.UnaryServerInterceptor))
	callPrice.RegisterHisDemoListServer(s, new(server.CallPriceServer))

	if err = s.Serve(l); err != nil {
		logrus.Fatal("failed to serve: %v", err)
	}
}
