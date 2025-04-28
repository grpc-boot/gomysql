package filter

import (
	"github.com/grpc-boot/gomysql/condition"
)

var (
	DefaultPageSize    = int64(20)
	DefaultPageSizeMax = int64(100)
)

type Page struct {
	Filters  map[string]Filter `json:"filters"`
	Sorts    []SortItem        `json:"sorts"`
	Page     int64             `json:"page"`
	PageSize int64             `json:"pageSize"`
}

func (p *Page) GetPage() int64 {
	if p.Page < 1 {
		p.Page = 1
	}
	return p.Page
}

func (p *Page) GetPageSize() int64 {
	if p.PageSize < 1 {
		p.PageSize = DefaultPageSize
	}

	if p.PageSize > DefaultPageSizeMax {
		p.PageSize = DefaultPageSizeMax
	}

	return p.PageSize
}

func (p *Page) Offset() int64 {
	return (p.GetPage() - 1) * p.GetPageSize()
}

func (p *Page) GetConditions() []condition.Condition {
	var conditions = make([]condition.Condition, 0)
	if len(p.Filters) == 0 {
		return conditions
	}

	for _, filter := range p.Filters {
		if cond := filter.GetCondition(); cond != nil {
			conditions = append(conditions, cond)
		}
	}

	return conditions
}

func (p *Page) GetOrders() []string {
	var orders = make([]string, 0)
	if len(p.Sorts) == 0 {
		return orders
	}

	for _, order := range p.Sorts {
		if orderStr := order.GetOrder(); orderStr != "" {
			orders = append(orders, orderStr)
		}
	}

	return orders
}
