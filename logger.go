package gomysql

var (
	logger   LogSql
	errorLog ErrLog
)

type (
	LogSql func(query string, args ...any)
	ErrLog func(err error, query string, args ...any)
)

func SetLogger(l LogSql) {
	logger = l
}

func SetErrorLog(l ErrLog) {
	errorLog = l
}
