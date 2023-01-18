package service

import (
	"context"
	"fmt"

	callPrice "github.com/CRORCR/duoo-common/proto/call_price"
)

type userService struct {
}

var UserService *userService

// CallPrice 获取主播私聊价格
func (s *userService) CallPrice(ctx context.Context, req *callPrice.GetPriceReq) (resp *callPrice.GetPriceResp, err error) {
	fmt.Println("收到请求：", req.Uid, ctx.Value("hello"))
	resp = new(callPrice.GetPriceResp)
	data := &callPrice.GetPriceResp_Data{Uid: "123", Date: "2023-01-18"}
	resp.Data = append(resp.Data, data)

	return resp, nil
}
