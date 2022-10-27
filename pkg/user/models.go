package user

import (
	"gorm.io/gorm"
	"time"
)

// User 用户信息
type User struct {
	gorm.Model
	Id       string
	Name     string
	Phone    string
	Gender   int
	ExpireAt time.Time
}

type Order struct {
	gorm.Model
	OrderId           string
	UserId            string
	Money             float64
	OrderStartTime    time.Time
	OrderFinishedTime time.Time
	Invalid           bool
}

// UserIP use to store user's device ip
type UserIP struct {
	gorm.Model
	UserId string
	Ip     string
	Desc   string
	Os     string
}

// UserMask mask for user
type UserMask struct {
	gorm.Model
	UserId string
	Mask   string
}

func (order Order) TableName() string {
	return "st_order"
}

func (u User) TableName() string {
	return "st_user"
}

func (ip UserIP) TableName() string {
	return "st_user_ip"
}

func (mask UserMask) TableName() string {
	return "st_user_mask"
}

// Create add a user
func (u User) Create(db *Db) error {
	return db.Create(&u).Error
}

// ListUser list users limit 10
func (u User) ListUser(db *Db) (users []User, err error) {
	err = db.Limit(10).Find(&users).Error
	return
}

func (u User) Get(db *Db) (user User, err error) {
	err = db.Where("id = ?", u.ID).Find(&user).Error
	return
}

func (u User) Delete(db *Db) error {
	return db.Delete(&u).Error
}

// Create add a user
func (mask UserMask) Create(db *Db) error {
	return db.Create(&mask).Error
}

// ListUser list users limit 10
func (mask UserMask) ListUser(db *Db) (userMasks []UserMask, err error) {
	err = db.Limit(10).Find(&userMasks).Error
	return
}

func (mask UserMask) Get(db *Db) (userMask UserMask, err error) {
	err = db.Where("id = ?", mask.ID).Find(&userMask).Error
	return
}

func (mask UserMask) Delete(db *Db) error {
	return db.Delete(&mask).Error
}

func (ip UserIP) Create(db *Db) error {
	return db.Create(&ip).Error
}

// ListUser list users limit 10
func (ip UserIP) ListUser(db *Db) (userIps []UserIP, err error) {
	err = db.Limit(10).Find(&userIps).Error
	return
}

func (ip UserIP) Get(db *Db) (userIp UserIP, err error) {
	err = db.Where("id = ?", ip.ID).Find(&userIp).Error
	return
}

func (ip UserIP) Delete(db *Db) error {
	return db.Delete(&ip).Error
}
