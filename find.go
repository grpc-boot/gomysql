package gomysql

import (
	"context"
	"time"

	"github.com/grpc-boot/gomysql/helper"
)

func FindModels[T Model](model T, db Executor, q *helper.Query) ([]T, error) {
	return FindModelsContext(context.Background(), model, db, q)
}

func FindModel[T Model](model T, db Executor, q *helper.Query) (T, error) {
	return FindModelContext(context.Background(), model, db, q)
}

func FindModelsTimeout[T Model](timeout time.Duration, model T, db Executor, q *helper.Query) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return FindModelsContext(ctx, model, db, q)
}

func FindModelTimeout[T Model](timeout time.Duration, model T, db Executor, q *helper.Query) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return FindModelContext(ctx, model, db, q)
}

func FindModelsContext[T Model](ctx context.Context, model T, db Executor, q *helper.Query) ([]T, error) {
	rows, err := SelectContext(ctx, db, q)
	if err != nil {
		return nil, err
	}

	return ScanModel(model, rows, err)
}

func FindModelContext[T Model](ctx context.Context, model T, db Executor, q *helper.Query) (m T, err error) {
	rows, err := SelectContext(ctx, db, q)
	if err != nil {
		return
	}

	models, err := ScanModel(model, rows, err)
	if err != nil {
		return
	}

	if len(models) == 0 {
		return
	}

	return models[0], nil
}

func FindOne(db Executor, q *helper.Query) (Record, error) {
	return FindOneContext(context.Background(), db, q)
}

func FindOneContext(ctx context.Context, db Executor, q *helper.Query) (Record, error) {
	records, err := FindContext(ctx, db, q)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records[0], nil
}

func FindOneTimeout(timeout time.Duration, db Executor, q *helper.Query) (Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return FindOneContext(ctx, db, q)
}

func Find(db Executor, q *helper.Query) (records []Record, err error) {
	rows, err := Select(db, q)
	return Scan(rows, err)
}

func FindContext(ctx context.Context, db Executor, q *helper.Query) (records []Record, err error) {
	rows, err := SelectContext(ctx, db, q)
	return Scan(rows, err)
}

func FindTimeout(timeout time.Duration, db Executor, q *helper.Query) (records []Record, err error) {
	rows, err := SelectTimeout(timeout, db, q)
	return Scan(rows, err)
}

func FindModelsByPoolContext[T Model](ctx context.Context, dbType DbType, model T, pool *Pool, q *helper.Query) ([]T, error) {
	query, args := q.Sql()
	rows, err := pool.QueryContext(ctx, dbType, query, args...)
	if err != nil {
		return nil, err
	}

	return ScanModel(model, rows, err)
}

func FindModelByPoolContext[T Model](ctx context.Context, dbType DbType, model T, pool *Pool, q *helper.Query) (m T, err error) {
	query, args := q.Sql()
	rows, err := pool.QueryContext(ctx, dbType, query, args...)
	if err != nil {
		return
	}

	models, err := ScanModel(model, rows, err)
	if err != nil {
		return
	}

	if len(models) == 0 {
		return
	}

	return models[0], nil
}

func FindModelsByPoolTimeout[T Model](timeout time.Duration, dbType DbType, model T, pool *Pool, q *helper.Query) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return FindModelsByPoolContext(ctx, dbType, model, pool, q)
}

func FindModelByPoolTimeout[T Model](timeout time.Duration, dbType DbType, model T, pool *Pool, q *helper.Query) (m T, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return FindModelByPoolContext(ctx, dbType, model, pool, q)
}
