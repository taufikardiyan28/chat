package UserRepo

import (
	interfaces "github.com/taufikardiyan28/chat/interfaces"
	UserModel "github.com/taufikardiyan28/chat/model/user"
)

type Repo struct {
	Pool interfaces.IDatabase
}

func (c *Repo) GetUserInfo(id string) (UserModel.User, error) {
	strSQL := `SELECT id, name, email, notifToken, IFNULL(lastSeen,'') AS lastSeen, status FROM users WHERE id=?`
	res := UserModel.User{}
	err := c.Pool.Get(&res, strSQL, id)
	return res, err
}

func (c *Repo) UpdateUser(id string, cols []string, val ...interface{}) error {
	strSQL := `UPDATE users SET `
	for _, col := range cols {
		strSQL += "`" + col + "` = ?"
	}
	strSQL += "WHERE `id` = ?"
	_, err := c.Pool.Exec(strSQL, id, val)
	return err
}
