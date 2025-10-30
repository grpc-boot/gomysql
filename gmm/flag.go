package gmm

import (
	"flag"
	"os"
)

var (
	DefaultFlag = &Flag{
		Host: "localhost",
		Port: 3306,
		User: "root",
	}
)

type Flag struct {
	Host      string
	Port      int
	Db        string
	User      string
	Passwd    string
	CharSet   string
	Table     string
	OutputDir string
}

func (f *Flag) Parse() {
	flag.StringVar(&f.Host, "h", "localhost", "Mysql host")
	flag.IntVar(&f.Port, "P", 3306, "Mysql port")
	flag.StringVar(&f.Db, "d", "", "Mysql db name")
	flag.StringVar(&f.User, "u", "root", "Mysql user")
	flag.StringVar(&f.Passwd, "p", "", "Mysql password")
	flag.StringVar(&f.CharSet, "c", "utf8", "Charset")
	flag.StringVar(&f.Table, "t", "", "Gmm table, All tables are selected by default")
	flag.StringVar(&f.OutputDir, "o", "entity", "Gmm model output directory")
	flag.Parse()

	if f.Db == "" {
		flag.Usage()
		os.Exit(1)
	}
}
