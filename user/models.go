package user

import "time"

// User 用户信息
type User struct {
	Id       string
	Name     string
	Phone    string
	Gender   int
	CreateAt time.Time
	DeleteAt time.Time
	ExpireAt time.Time
}

type Order struct {
	OrderId           string
	UserId            string
	Money             float64
	OrderStartTime    time.Time
	OrderFinishedTime time.Time
	Invalid           bool
}

type Conf struct {
	Listen string `mapstruce:"listen"`
	Mysql  struct {
		Url      string
		User     string
		Password string
	} `mapsruce:"mysql"`
}

func (order Order) TableName() string {
	return "star_order"
}

func (user User) TableName() string {
	return "star_user"
}

func (u *User) Create() error {
	return nil
}
