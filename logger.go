package gomysql

var (
	logger LogSql
)

type LogSql func(query string, args ...any)

func SetLogger(l LogSql) {
	logger = l
}
