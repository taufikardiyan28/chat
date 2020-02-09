package UserModel

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"-"`
	NotifToken string `json:"-" db:"notifToken"`
	LastSeen   string `json:"last_seen" db:"lastSeen"`
	Status     string `json:"status" db:"status"`
}
