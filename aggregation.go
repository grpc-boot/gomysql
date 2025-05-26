package gomysql

import (
	"context"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func CountByConditionTimeout[T DbModel](timeout time.Duration, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return CountByConditionContext(ctx, model, db, cond, field, tableArgs...)
}

func CountByConditionContext[T DbModel](ctx context.Context, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (count int64, err error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(cond)
	)

	defer q.Close()
	return CountContext(ctx, db, q, field)
}

func CountTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return CountContext(ctx, db, q, field)
}

func CountContext(ctx context.Context, db Executor, q *helper.Query, field string) (count int64, err error) {
	query, args := q.CountSql(field)
	rows, err := QueryContext(ctx, db, query, args...)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].ToInt64("countVal"), nil
}

func SumByConditionTimeout[T DbModel](timeout time.Duration, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (sum float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return SumByConditionContext(ctx, model, db, cond, field, tableArgs...)
}

func SumByConditionContext[T DbModel](ctx context.Context, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (sum float64, err error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(cond)
	)

	defer q.Close()
	return SumContext(ctx, db, q, field)
}

func SumTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (sum float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return SumContext(ctx, db, q, field)
}

func SumContext(ctx context.Context, db Executor, q *helper.Query, field string) (sum float64, err error) {
	query, args := q.SumSql(field)
	rows, err := QueryContext(ctx, db, query, args...)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].ToFloat64("sumVal"), nil
}

func MaxByConditionTimeout[T DbModel](timeout time.Duration, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (max string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return MaxByConditionContext(ctx, model, db, cond, field, tableArgs...)
}

func MaxByConditionContext[T DbModel](ctx context.Context, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (max string, err error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(cond)
	)

	defer q.Close()
	return MaxContext(ctx, db, q, field)
}

func MaxTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (max string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return MaxContext(ctx, db, q, field)
}

func MaxContext(ctx context.Context, db Executor, q *helper.Query, field string) (max string, err error) {
	query, args := q.MaxSql(field)
	rows, err := QueryContext(ctx, db, query, args...)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].String("maxVal"), nil
}

func MinByConditionTimeout[T DbModel](timeout time.Duration, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (min string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return MinByConditionContext(ctx, model, db, cond, field, tableArgs...)
}

func MinByConditionContext[T DbModel](ctx context.Context, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (min string, err error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(cond)
	)

	defer q.Close()
	return MinContext(ctx, db, q, field)
}

func MinTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (min string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return MinContext(ctx, db, q, field)
}

func MinContext(ctx context.Context, db Executor, q *helper.Query, field string) (min string, err error) {
	query, args := q.MinSql(field)
	rows, err := QueryContext(ctx, db, query, args...)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].String("minVal"), nil
}

func AvgByConditionTimeout[T DbModel](timeout time.Duration, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (avg string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return AvgByConditionContext(ctx, model, db, cond, field, tableArgs...)
}

func AvgByConditionContext[T DbModel](ctx context.Context, model T, db Executor, cond condition.Condition, field string, tableArgs ...any) (avg string, err error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(cond)
	)

	defer q.Close()
	return AvgContext(ctx, db, q, field)
}

func AvgTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (avg string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return AvgContext(ctx, db, q, field)
}

func AvgContext(ctx context.Context, db Executor, q *helper.Query, field string) (avg string, err error) {
	query, args := q.AvgSql(field)
	rows, err := QueryContext(ctx, db, query, args...)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].String("avgVal"), nil
}
