package MySqlDB

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/taufikardiyan28/chat/helper"
)

type Conn struct {
	Config *helper.Configuration
	Pool   *sqlx.DB
}

func (c *Conn) Connect() error {
	cfg := c.Config
	dbConfig := mysql.Config{
		User:                 cfg.Database.User,
		Passwd:               cfg.Database.Password,
		DBName:               cfg.Database.DbName,
		Loc:                  time.Local,
		Net:                  fmt.Sprintf("tcp(%s:%d)", cfg.Database.Host, cfg.Database.Port),
		AllowNativePasswords: true,
		MultiStatements:      true,
		ParseTime:            true,
	}

	dsn := dbConfig.FormatDSN()
	var err error

	c.Pool, err = sqlx.Open("mysql", dsn)
	c.Pool.SetMaxIdleConns(50)
	c.Pool.SetMaxOpenConns(50)
	return err
}

func (c *Conn) Ping() error {
	return c.Pool.Ping()
}

func (c *Conn) GetPool() interface{} {
	return c.Pool
}
