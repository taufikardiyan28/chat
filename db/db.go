package db

import (
	MongoDB "github.com/taufikardiyan28/chat/db/mongo"
	MySqlDB "github.com/taufikardiyan28/chat/db/mysql"
	"github.com/taufikardiyan28/chat/helper"
	"github.com/taufikardiyan28/chat/interfaces"
)

func NewConnection(config *helper.Configuration) (interfaces.Database, error) {
	var iDB interfaces.Database
	if config.Database.DbType == "mongodb" {
		iDB = &MongoDB.Conn{
			Config: config,
		}
	} else {
		iDB = &MySqlDB.Conn{
			Config: config,
		}
	}

	err := iDB.Connect()
	return iDB, err
}
