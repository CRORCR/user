package service

import (
	"math"

	"github.com/CRORCR/cr-common/proto/base"
)

var BaseService = &baseService{}

type baseService struct {
}

// 分页处理
func (*baseService) NewPagination(total, size, page int32) *base.Pagination {
	totalPage := int32(math.Ceil(float64(total) / float64(size)))
	return &base.Pagination{
		TotalPages:  totalPage,
		TotalRecord: total,
		CurrentPage: page,
		PageSize:    size,
		HasNext:     totalPage > page,
	}
}
