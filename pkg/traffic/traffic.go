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
// limitations under the License.

package traffic

import (
	"sync"
	"time"
)

const (
	normalSize = 1024 * 1024 * 1024 * 5 //2GB

)

type traffic struct {
	lock  sync.Mutex
	Size  int64
	start time.Time
	end   time.Time
}

type Interface interface {
	Add(userId string, size int64) error
	Clear(userId string) error
}

func NewTraffic() Interface {
	return &traffic{}
}

func (t *traffic) Add(userId string, size int64) error {
	return nil
}

func (t *traffic) Clear(userId string) error {
	return nil
}
