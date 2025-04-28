package filter

import "strings"

type SortItem struct {
	Field string `json:"field"`
	Kind  string `json:"kind"`
}

func (si *SortItem) IsDesc() bool {
	return strings.EqualFold(si.Kind, "desc")
}

func (si *SortItem) GetOrder() string {
	if si.Field == "" {
		return ""
	}

	if si.IsDesc() {
		return si.Field + " DESC"
	}
	return si.Field + " ASC"
}
