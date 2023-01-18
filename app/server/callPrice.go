package server

import (
	"context"

	callPrice "github.com/CRORCR/duoo-common/proto/call_price"
	"github.com/CRORCR/user/internal/service"
)

// CallPriceServer 聊价
type CallPriceServer struct {
}

func (c *CallPriceServer) GetDemo(ctx context.Context, req *callPrice.GetPriceReq) (resp *callPrice.GetPriceResp, err error) {
	return service.UserService.CallPrice(ctx, req)
}
