package util

import "sync"

var (
	UCTL = user{
		UserId: "123456789abcdef0",
	}
)

type UserFunc interface {
	GetUserId() ([]byte, error)
	SetUserId(userId string) error
	SetUserInfo(username, password string) error
}

// User user username password to login, then will receive userId
type user struct {
	lock     sync.Mutex
	Username string
	Password string
	UserId   string
}

func NewUser() UserFunc {
	return &user{}
}

var (
	_ UserFunc = (*user)(nil)
)

func (u *user) GetUserId() ([]byte, error) {
	return nil, nil
}

func (u *user) SetUserId(userId string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	UCTL.UserId = userId
	return nil
}

func (u *user) SetUserInfo(username, password string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	UCTL.Username = username
	UCTL.Password = password
	return nil
}
