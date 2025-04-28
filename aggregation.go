package gomysql

import (
	"context"
	"time"

	"github.com/grpc-boot/gomysql/helper"
)

func CountTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return CountContext(ctx, db, q, field)
}

func CountContext(ctx context.Context, db Executor, q *helper.Query, field string) (count int64, err error) {
	query, args := q.CountSql(field)
	rows, err := QueryContext(ctx, db, query, args)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].ToInt64("countVal"), nil
}

func SumTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (sum float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return SumContext(ctx, db, q, field)
}

func SumContext(ctx context.Context, db Executor, q *helper.Query, field string) (sum float64, err error) {
	query, args := q.SumSql(field)
	rows, err := QueryContext(ctx, db, query, args)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].ToFloat64("sumVal"), nil
}

func MaxTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (max string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return MaxContext(ctx, db, q, field)
}

func MaxContext(ctx context.Context, db Executor, q *helper.Query, field string) (max string, err error) {
	query, args := q.MaxSql(field)
	rows, err := QueryContext(ctx, db, query, args)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].String("maxVal"), nil
}

func MinTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (min string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return MinContext(ctx, db, q, field)
}

func MinContext(ctx context.Context, db Executor, q *helper.Query, field string) (min string, err error) {
	query, args := q.MinSql(field)
	rows, err := QueryContext(ctx, db, query, args)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].String("minVal"), nil
}

func AvgTimeout(timeout time.Duration, db Executor, q *helper.Query, field string) (avg string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return AvgContext(ctx, db, q, field)
}

func AvgContext(ctx context.Context, db Executor, q *helper.Query, field string) (avg string, err error) {
	query, args := q.AvgSql(field)
	rows, err := QueryContext(ctx, db, query, args)
	if err != nil {
		return
	}

	records, err := Scan(rows, err)
	if err != nil || len(records) == 0 {
		return
	}

	return records[0].String("avgVal"), nil
}
