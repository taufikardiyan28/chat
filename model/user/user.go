package UserModel

type User struct {
	ID         string `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Email      string `json:"email" db:"email"`
	Phone      string `json:"phone" db:"phone"`
	NotifToken string `json:"-" db:"notifToken"`
	LastSeen   string `json:"last_seen" db:"lastSeen"`
	Status     string `json:"status" db:"status"`
	IsActive   int64  `json:"is_active" db:"isActive"`
}
