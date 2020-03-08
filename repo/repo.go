package Repo

import (
	"github.com/taufikardiyan28/chat/interfaces"
	MysqlMessageRepo "github.com/taufikardiyan28/chat/repo/mysql/messages"
	MysqlUserRepo "github.com/taufikardiyan28/chat/repo/mysql/user"
)

func GetUserRepo(pool interfaces.IDatabase) interfaces.IUserRepo {
	var iFace interfaces.IUserRepo
	/*if config.Database.DbType == "mongodb" {
		iFace = &MongoUserRepo.User{
			Pool: pool,
		}
	} else {
		iFace = &MysqlUserRepo.User{
			Pool: pool,
		}
	}*/

	iFace = &MysqlUserRepo.Repo{
		Pool: pool,
	}
	return iFace
}

func GetMessageRepo(pool interfaces.IDatabase) interfaces.IMessageRepo {
	var iFace interfaces.IMessageRepo
	/*if config.Database.DbType == "mongodb" {
		iFace = &MysqlMessageRepo.User{
			Pool: pool,
		}
	} else {
		iFace = &MysqlMessageRepo.User{
			Pool: pool,
		}
	}*/

	iFace = &MysqlMessageRepo.Repo{
		Pool: pool,
	}
	return iFace
}
