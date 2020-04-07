package db

import (
	"com/database/orm"
	"database/sql"
	"fmt"
	"testing"
)

var cfg = &orm.DbCfg{
	Driver:  "mysql",
	Usr:     "root",
	Pswd:    "123",
	Host:    "127.0.0.1",
	Port:    3306,
	Db:      "",
	Charset: "utf8",
}

func TestDB_Init(t *testing.T) {
	conStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		cfg.Usr, cfg.Pswd, cfg.Host, cfg.Port, cfg.Db, cfg.Charset)
	c, err := sql.Open(cfg.Driver, conStr)
	if err != nil {
		t.Error("conn db err:v", err)
	}
	r, err := c.Exec("CREATE DATABASE  IF NOT EXISTS game")
	if err != nil {
		t.Error(err)
	}
	r, err = c.Exec("use game")
	if err != nil {
		t.Error(err)
	}
	r, err = c.Exec("show tables")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)
}
