package filter

import (
	"strconv"

	"github.com/grpc-boot/gomysql/condition"
)

type Scroll struct {
	Filters  map[string]Filter `json:"filters"`
	Sorts    SortItem          `json:"sorts"`
	Cursor   string            `json:"cursor"`
	PageSize int64             `json:"pageSize"`
}

func (s *Scroll) CursorInt(defaultVal int64) (cursor int64) {
	cursor, err := strconv.ParseInt(s.Cursor, 10, 64)
	if err != nil {
		return defaultVal
	}
	return cursor
}

func (s *Scroll) GetPageSize() int64 {
	if s.PageSize < 1 {
		s.PageSize = DefaultPageSize
	}

	if s.PageSize > DefaultPageSizeMax {
		s.PageSize = DefaultPageSizeMax
	}

	return s.PageSize
}

func (s *Scroll) GetConditions() []condition.Condition {
	var conditions = make([]condition.Condition, 0)
	if len(s.Filters) == 0 {
		return conditions
	}

	for _, filter := range s.Filters {
		if cond := filter.GetCondition(); cond != nil {
			conditions = append(conditions, cond)
		}
	}

	return conditions
}

func (s *Scroll) GetOrders() []string {
	var orders = make([]string, 0)
	if orderStr := s.Sorts.GetOrder(); orderStr != "" {
		orders = append(orders, orderStr)
	}
	return orders
}
