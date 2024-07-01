package gomysql

import (
	"context"
	"time"

	"github.com/grpc-boot/gomysql/helper"
)

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
