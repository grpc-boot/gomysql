package gomysql

import (
	"database/sql"

	"github.com/grpc-boot/gomysql/convert"
)

func ScanStringValues(rows *sql.Rows, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err = rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

func Scan(rows *sql.Rows, err error) ([]Record, error) {
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var (
		values  = make([]any, len(fields), len(fields))
		records = make([]Record, 0, 8)
	)

	for index, _ := range fields {
		values[index] = &[]byte{}
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		record := make(map[string]string, len(fields))
		for index, field := range fields {
			record[field] = convert.Bytes2String(*values[index].(*[]byte))
		}
		records = append(records, record)
	}

	return records, nil
}

func ScanBytesRecords(rows *sql.Rows, err error) ([]BytesRecord, error) {
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fields, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var (
		values  = make([]any, len(fields), len(fields))
		records = make([]BytesRecord, 0, 8)
	)

	for index, _ := range fields {
		values[index] = &[]byte{}
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		record := make(BytesRecord, len(fields))
		for index, field := range fields {
			record[field] = *values[index].(*[]byte)
		}
		records = append(records, record)
	}

	return records, nil
}

func ScanModel[T Model](model T, rows *sql.Rows, err error) ([]T, error) {
	brs, err := ScanBytesRecords(rows, err)
	if err != nil {
		return nil, err
	}

	return BytesRecords2Models(brs, model), nil
}
