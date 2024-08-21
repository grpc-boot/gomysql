package gomysql

import (
	"context"
	"time"

	"github.com/grpc-boot/gomysql/condition"
	"github.com/grpc-boot/gomysql/helper"
)

type DbModel interface {
	Model
	PrimaryKey() string
	TableName(args ...any) string
}

func FindById[T DbModel](db Executor, id int64, model T, args ...any) (T, error) {
	q := helper.AcquireQuery().
		From(model.TableName(args...)).
		Where(condition.Equal{Field: model.PrimaryKey(), Value: id}).
		Limit(1)
	defer q.Close()

	return FindModel[T](model, db, q)
}

func FindAll[T DbModel](db Executor, con condition.Condition, model T, args ...any) ([]T, error) {
	q := helper.AcquireQuery().
		From(model.TableName(args...)).
		Where(con)
	defer q.Close()

	return FindModels[T](model, db, q)
}

func FindByIdWithPoolContext[T DbModel](ctx context.Context, pool *Pool, dbType DbType, id int64, model T, args ...any) (T, error) {
	q := helper.AcquireQuery().
		From(model.TableName(args...)).
		Where(condition.Equal{Field: model.PrimaryKey(), Value: id}).
		Limit(1)
	defer q.Close()

	return FindModelByPoolContext(ctx, dbType, model, pool, q)
}

func FindAllWithPoolContext[T DbModel](ctx context.Context, pool *Pool, dbType DbType, con condition.Condition, model T, args ...any) ([]T, error) {
	q := helper.AcquireQuery().
		From(model.TableName(args...)).
		Where(con)
	defer q.Close()

	return FindModelsByPoolContext(ctx, dbType, model, pool, q)
}

func FindByIdTimeoutWithPool[T DbModel](timeout time.Duration, pool *Pool, dbType DbType, id int64, model T, args ...any) (T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return FindByIdWithPoolContext(ctx, pool, dbType, id, model, args...)
}

func FindAllTimeoutWithPool[T DbModel](timeout time.Duration, pool *Pool, dbType DbType, con condition.Condition, model T, args ...any) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return FindAllWithPoolContext(ctx, pool, dbType, con, model, args...)
}
