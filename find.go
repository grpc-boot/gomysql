package gomysql

import (
	"context"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

func FindModelByIdTimeout[T DbModel](timeout time.Duration, model T, db Executor, id any, tableArgs ...any) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return FindModelByIdContext(ctx, model, db, id, tableArgs...)
}

func FindModelByIdContext[T DbModel](ctx context.Context, model T, db Executor, id any, tableArgs ...any) (T, error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(condition.Equal{Field: model.PrimaryKey(), Value: id})
	)

	defer q.Close()
	return FindModelContext(ctx, model, db, q)
}

func FindModelByConditionTimeout[T DbModel](timeout time.Duration, model T, db Executor, cond condition.Condition, tableArgs ...any) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return FindModelByConditionContext(ctx, model, db, cond, tableArgs...)
}

func FindModelByConditionContext[T DbModel](ctx context.Context, model T, db Executor, cond condition.Condition, tableArgs ...any) (T, error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(cond)
	)

	defer q.Close()
	return FindModelContext(ctx, model, db, q)
}

func FindModelsByConditionTimeout[T DbModel](timeout time.Duration, model T, db Executor, cond condition.Condition, tableArgs ...any) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return FindModelsByConditionContext(ctx, model, db, cond, tableArgs...)
}

func FindModelsByConditionContext[T DbModel](ctx context.Context, model T, db Executor, cond condition.Condition, tableArgs ...any) ([]T, error) {
	var (
		q = helper.AcquireQuery().
			From(model.TableName(tableArgs...)).
			Where(cond)
	)

	defer q.Close()
	return FindModelsContext(ctx, model, db, q)
}

func FindModels[T Model](model T, db Executor, q *helper.Query) ([]T, error) {
	return FindModelsContext(context.Background(), model, db, q)
}

func FindModelsTimeout[T Model](timeout time.Duration, model T, db Executor, q *helper.Query) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return FindModelsContext(ctx, model, db, q)
}

func FindModelsContext[T Model](ctx context.Context, model T, db Executor, q *helper.Query) ([]T, error) {
	rows, err := SelectContext(ctx, db, q)
	if err != nil {
		return make([]T, 0), err
	}

	return ScanModel(model, rows, err)
}

func FindModel[T Model](model T, db Executor, q *helper.Query) (T, error) {
	return FindModelContext(context.Background(), model, db, q)
}

func FindModelTimeout[T Model](timeout time.Duration, model T, db Executor, q *helper.Query) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return FindModelContext(ctx, model, db, q)
}

func FindModelContext[T Model](ctx context.Context, model T, db Executor, q *helper.Query) (m T, err error) {
	models, err := FindModelsContext(ctx, model, db, q)
	if err != nil {
		return
	}

	if len(models) > 0 {
		return models[0], nil
	}
	return
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
	return FindContext(context.Background(), db, q)
}

func FindContext(ctx context.Context, db Executor, q *helper.Query) (records []Record, err error) {
	rows, err := SelectContext(ctx, db, q)
	return Scan(rows, err)
}

func FindTimeout(timeout time.Duration, db Executor, q *helper.Query) (records []Record, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return FindContext(ctx, db, q)
}
