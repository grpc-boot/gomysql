package gomysql

import "fmt"

func ShowCreateTable(db Executor, table string) (sqlStr string, err error) {
	var name string
	err = QueryRow(db, fmt.Sprintf("SHOW CREATE TABLE `%s`", table)).Scan(&name, &sqlStr)
	return sqlStr, err
}

func ShowTables(db Executor) ([]string, error) {
	rows, err := Query(db, "SHOW TABLES")
	if err != nil {
		return nil, err
	}

	return ScanStringValues(rows, err)
}
