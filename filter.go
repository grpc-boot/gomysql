package gomysql

import (
	"context"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/filter"
	"github.com/grpc-boot/gomysql/helper"
)

func ScrollFilterWithIdTimeout[T DbModel](timeout time.Duration, filter *filter.Scroll, defaultId int64, model T, db Executor, args ...any) (hasMore bool, last T, ms []T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ScrollFilterWithIdContext(ctx, filter, defaultId, model, db, args...)
}

func ScrollFilterWithIdContext[T DbModel](ctx context.Context, filter *filter.Scroll, defaultId int64, model T, db Executor, args ...any) (hasMore bool, last T, ms []T, err error) {
	var (
		q          = helper.AcquireQuery()
		where      = condition.And(filter.GetConditions())
		orders     = filter.GetOrders()
		fetchSize  = filter.GetPageSize() + 1
		offset     = filter.CursorInt(defaultId)
		scrollCond = condition.Condition(condition.Lte{Field: model.PrimaryKey(), Value: offset})
	)

	defer q.Close()

	q.From(model.TableName(args...)).
		Limit(fetchSize)

	if len(orders) > 0 {
		q.Order(orders...)
		if !filter.Sorts.IsDesc() {
			scrollCond = condition.Gte{Field: model.PrimaryKey(), Value: offset}
		}
	} else {
		q.Order(model.PrimaryKey() + " DESC")
	}

	where = append(where, scrollCond)
	q.Where(where)

	ms, err = FindModelsContext(ctx, model, db, q)
	if len(ms) > 0 {
		last = ms[len(ms)-1]
	}

	if len(ms) > int(filter.GetPageSize()) {
		hasMore = true
		ms = ms[:len(ms)-1]
	}
	return
}

func PageFilterWithTotalTimeout[T DbModel](timeout time.Duration, filter *filter.Page, model T, db Executor, args ...any) (totalCount, pageCount int64, ms []T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return PageFilterWithTotalContext(ctx, filter, model, db, args...)
}

func PageFilterWithTotalContext[T DbModel](ctx context.Context, filter *filter.Page, model T, db Executor, args ...any) (totalCount, pageCount int64, ms []T, err error) {
	var (
		q      = helper.AcquireQuery()
		where  = condition.And(filter.GetConditions())
		orders = filter.GetOrders()
	)

	defer q.Close()

	q.From(model.TableName(args...))
	if len(where) > 0 {
		q.Where(where)
	}

	totalCount, err = CountContext(ctx, db, q, "*")
	if err != nil || totalCount < 1 {
		return
	}

	pageCount = helper.PageCount(totalCount, filter.GetPageSize())
	if filter.GetPage() > pageCount {
		return
	}

	if len(orders) > 0 {
		q.Order(orders...)
	}

	q.Offset(filter.Offset()).
		Limit(filter.GetPageSize())
	ms, err = FindModelsContext(ctx, model, db, q)
	return
}

func PageFilterTimeout[T DbModel](timeout time.Duration, filter *filter.Page, model T, db Executor, args ...any) (ms []T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return PageFilterContext(ctx, filter, model, db, args...)
}

func PageFilterContext[T DbModel](ctx context.Context, filter *filter.Page, model T, db Executor, args ...any) (ms []T, err error) {
	var (
		q      = helper.AcquireQuery()
		where  = filter.GetConditions()
		orders = filter.GetOrders()
	)

	defer q.Close()

	q.From(model.TableName(args...)).
		Offset(filter.Offset()).
		Limit(filter.GetPageSize())
	if len(where) > 0 {
		q.Where(condition.And(where))
	}

	if len(orders) > 0 {
		q.Order(orders...)
	}

	return FindModelsContext(ctx, model, db, q)
}
