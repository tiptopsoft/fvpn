// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import "sync"

type UserFunc interface {
	GetUserId() string
	SetUserId(userId string) error
	SetUserInfo(username, password string) error
}

// User username password to login, then will receive userId
type user struct {
	lock     sync.Mutex
	username string
	password string
	userId   string
}

func Info() UserFunc {
	return &user{}
}

var (
	_ UserFunc = (*user)(nil)
)

func (u *user) GetUserId() string {
	return u.userId
}

func (u *user) SetUserId(userId string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.userId = userId
	return nil
}

func (u *user) SetUserInfo(username, password string) error {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.username = username
	u.password = password
	return nil
}
