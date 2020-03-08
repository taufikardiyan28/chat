package MySqlDB

import UserModel "github.com/taufikardiyan28/chat/model/user"

func (c *Conn) GetUserInfo(id string) (UserModel.User, error) {
	strSQL := `SELECT id, name, email, notifToken, IFNULL(lastSeen,'') AS lastSeen, status FROM users WHERE id=?`
	res := UserModel.User{}
	err := c.Pool.Get(&res, strSQL, id)
	return res, err
}

func (c *Conn) UpdateUser(id string) error {
	strSQL := `UPDATE users SET name=?, email=?, notifToken=?, lastSeen=?, status=? WHERE id=?`
	_, err := c.Pool.Exec(strSQL, id)
	return err
}
