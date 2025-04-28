package filter

import (
	"strings"

	"github.com/grpc-boot/gomysql/condition"
)

type Filter struct {
	Field  string `json:"field"`
	Option string `json:"opt"`
	Value  string `json:"value"`
}

func (f *Filter) GetCondition() condition.Condition {
	if f.Field == "" || f.Option == "" {
		return nil
	}

	switch f.Option {
	case "=":
		if f.Value == "" {
			return nil
		}
		return condition.Equal{Field: f.Field, Value: f.Value}
	case "≠":
		if f.Value == "" {
			return nil
		}
		return condition.NotEqual{Field: f.Field, Value: f.Value}
	case ">":
		if f.Value == "" {
			return nil
		}
		return condition.Gt{Field: f.Field, Value: f.Value}
	case ">=":
		if f.Value == "" {
			return nil
		}
		return condition.Gte{Field: f.Field, Value: f.Value}
	case "<":
		if f.Value == "" {
			return nil
		}
		return condition.Lt{Field: f.Field, Value: f.Value}
	case "<=":
		if f.Value == "" {
			return nil
		}
		return condition.Lte{Field: f.Field, Value: f.Value}
	case "Contains":
		if f.Value == "" {
			return nil
		}
		return condition.Contains{Field: f.Field, Words: f.Value}
	case "Not Contains":
		if f.Value == "" {
			return nil
		}
		return condition.NotContains{Field: f.Field, Words: f.Value}
	case "Start With":
		if f.Value == "" {
			return nil
		}
		return condition.BeginWith{Field: f.Field, Words: f.Value}
	case "End With":
		if f.Value == "" {
			return nil
		}
		return condition.EndWith{Field: f.Field, Words: f.Value}
	case "Is Not Empty":
		return condition.NotEmpty{Field: f.Field}
	case "Is Empty":
		return condition.Empty{Field: f.Field}
	case "Is Any Of":
		return condition.In[string]{Field: f.Field, Value: strings.Split(f.Value, ",")}
	default:
		return nil
	}
}
