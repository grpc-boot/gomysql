package gomysql

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grpc-boot/gomysql/gmm"
	"github.com/grpc-boot/gomysql/inner"
)

func GmmConsole() {
	gmm.DefaultFlag.Parse()

	var (
		tables  []string
		db, err = NewDb(Options{
			DbName:   gmm.DefaultFlag.Db,
			Host:     gmm.DefaultFlag.Host,
			Port:     uint32(gmm.DefaultFlag.Port),
			UserName: gmm.DefaultFlag.User,
			Password: gmm.DefaultFlag.Passwd,
			CharSet:  gmm.DefaultFlag.CharSet,
		})
	)

	if err != nil {
		inner.RedFatal("connect db err: %v", err)
	}

	executor := db.Executor()

	if gmm.DefaultFlag.Table != "" {
		tables = []string{gmm.DefaultFlag.Table}
	} else {
		tables, err = ShowTables(executor)
		if err != nil {
			inner.RedFatal("show tables err: %v", err)
		}

		if len(tables) == 0 {
			inner.RedFatal("show tables err: table is empty")
		}
	}

	for _, table := range tables {
		var createSql string
		createSql, err = ShowCreateTable(executor, table)
		if err != nil {
			inner.RedFatal("show create table err: %v", err)
		}

		err = Gmm(gmm.DefaultFlag.OutputDir, table, createSql)
		if err != nil {
			inner.Red("create model for table:%s err: %v", table, err)
		} else {
			inner.Green("create model for table:%s success", table)
		}
	}
}

func GmmAll(dirPath string, db Executor, errContinue bool) error {
	tables, err := ShowTables(db)
	if err != nil {
		return err
	}

	if len(tables) == 0 {
		return nil
	}

	for _, table := range tables {
		var createSql string
		createSql, err = ShowCreateTable(db, table)
		if err != nil {
			return err
		}

		if err = Gmm(dirPath, table, createSql); err != nil {
			if !errContinue {
				return err
			}
		}
	}

	return nil
}

func Gmm(dirPath, tableName, tableCreateSql string) error {
	var pkg = filepath.Base(dirPath)
	if pkg == "." || pkg == "./" || pkg == "" || pkg == "/" || pkg == "../" {
		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			return err
		}

		pkg = filepath.Base(absPath)
	}

	if pkg == "." || pkg == "./" || pkg == "" || pkg == "/" || pkg == "../" {
		return errors.New("parse pkg name failed")
	}

	var (
		modelFileName = fmt.Sprintf("%s/%s.go", strings.TrimSuffix(dirPath, string(os.PathSeparator)), tableName)
		crudFileName  = fmt.Sprintf("%s/%s_crud.go", strings.TrimSuffix(dirPath, string(os.PathSeparator)), tableName)
		exists, err   = inner.FileExists(modelFileName)
	)

	if err != nil {
		return err
	}

	if exists {
		return ErrModelFileExists
	}

	exists, err = inner.FileExists(crudFileName)
	if err != nil {
		return err
	}

	if exists {
		return ErrCrudFileExists
	}

	exists, err = inner.FileExists(dirPath)
	if err != nil {
		return err
	}

	if !exists {
		if err = inner.MkDir(dirPath, 0755); err != nil {
			return err
		}
	}

	var (
		primaryKey          = "id"
		fields, primaryKeys = gmm.ParseCreateTable(tableCreateSql)
	)

	if len(primaryKeys) > 0 {
		primaryKey = primaryKeys[0]
	}

	model, crud := gmm.GenerateStruct(primaryKey, pkg, tableName, fields)
	if err = os.WriteFile(modelFileName, model, 0644); err != nil {
		return err
	}

	return os.WriteFile(crudFileName, crud, 0644)
}
